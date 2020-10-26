package scrape

import (
	"context"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"golang.org/x/time/rate"
	"google.golang.org/api/sheets/v4"
)

type sheetsSvc struct {
	// Rate limited service
	*rate.Limiter
	svc *sheets.Service
}

// Please don't hack me
const sheetID = "1wB-bJrtygi3Tab8n9XZgwINjHIqG-K4O_2j0QJggP18"
const dumpRange = "Dump!A2"

func newSheets(ctx context.Context) *sheetsSvc {
	client, err := google.DefaultClient(ctx, sheets.SpreadsheetsScope)
	limiter := rate.NewLimiter(1, 1)

	sheetsService, err := sheets.New(client)
	if err != nil {
		logrus.Fatalf("Failed to create sheets client. Err: %v", err)
	}

	return &sheetsSvc{limiter, sheetsService}
}

func (svc *sheetsSvc) storeRealestate(re realestate) {
	if !svc.Allow() {
		logrus.Warning("Waiting for Google API limit")
	}

	if err := svc.Wait(context.Background()); err != nil {
		logrus.Errorf("Failed to wait for rate limiter. Error: %v", err)
	}

	logrus.Debug("Storing [%s]: %s", re.ID, re.Title)
	vr := sheets.ValueRange{}
	vr.Values = append(vr.Values, re.toRow())
	_, err := svc.svc.Spreadsheets.Values.Append(sheetID, dumpRange, &vr).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		logrus.Errorf("Failed to write! Err: %v", err)
	}
}
