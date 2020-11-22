package sheet

import (
	"context"
	"os"

	"github.com/gmenih341/scrappr/src/realestate"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"golang.org/x/time/rate"
	"google.golang.org/api/sheets/v4"
)

type sheetService struct {
	*rate.Limiter
	*sheets.Service
	sheetID string
}

const dumpRange = "Dump!A2"

func NewService(ctx context.Context) *sheetService {
	client, err := google.DefaultClient(ctx, sheets.SpreadsheetsScope)
	limiter := rate.NewLimiter(1, 1)

	sheetsService, err := sheets.New(client)
	if err != nil {
		logrus.Fatalf("Failed to create sheets client. Err: %v", err)
	}

	return &sheetService{
		limiter,
		sheetsService,
		os.Getenv("SPREADSHEET_ID"),
	}
}

func (svc *sheetService) StoreRealestate(re realestate.RealestateEntity) {
	if !svc.Allow() {
		logrus.Warning("Waiting for Google API limit")
	}

	if err := svc.Wait(context.Background()); err != nil {
		logrus.Errorf("Failed to wait for rate limiter. Error: %v", err)
	}

	logrus.Debugf("Storing [%s]: %s", re.ID, re.Title)
	vr := sheets.ValueRange{}
	vr.Values = append(vr.Values, re.ToRow())
	_, err := svc.Spreadsheets.Values.Append(svc.sheetID, dumpRange, &vr).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		logrus.Errorf("Failed to write! Err: %v", err)
	}
}
