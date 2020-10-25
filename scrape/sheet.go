package scrape

import (
	"context"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"golang.org/x/time/rate"
	"google.golang.org/api/sheets/v4"
)

type sheetsOp struct {
	svc *sheets.Service
	lim *rate.Limiter
}

const sheetID = "1wB-bJrtygi3Tab8n9XZgwINjHIqG-K4O_2j0QJggP18"
const dumpRange = "Dump!A2"

func newOp(ctx context.Context) *sheetsOp {
	client, err := google.DefaultClient(ctx, sheets.SpreadsheetsScope)
	limiter := rate.NewLimiter(1, 1)

	sheetsService, err := sheets.New(client)
	if err != nil {
		logrus.Fatalf("Failed to create sheets client. Err: %v", err)
	}

	return &sheetsOp{sheetsService, limiter}
}

func (so *sheetsOp) storeApartment(ap house) {
	if !so.lim.Allow() {
		logrus.Warning("Waiting for Google API limit")
	}

	if err := so.lim.Wait(context.Background()); err != nil {
		logrus.Errorf("Failed to wait for rate limiter. Error: %v", err)
	}

	logrus.Debug("Storing apartment")
	vr := sheets.ValueRange{}
	vr.Values = append(vr.Values, ap.toRow())
	_, err := so.svc.Spreadsheets.Values.Append(sheetID, dumpRange, &vr).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		logrus.Errorf("Failed to write! Err: %v", err)
	}
}
