package scrape

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type storeOp struct {
	client *datastore.Client
	lim    *rate.Limiter
	ctx    context.Context
}

func newStore(ctx context.Context) *storeOp {
	dsClient, err := datastore.NewClient(ctx, "website-246917")
	limiter := rate.NewLimiter(1, 1)

	if err != nil {
		logrus.Fatalf("Failed to create sheets client. Err: %v", err)
	}

	return &storeOp{dsClient, limiter, ctx}
}

func (so *storeOp) storeApartment(ap house) {
	if !so.lim.Allow() {
		logrus.Warning("Waiting for Google API limit")
	}

	if err := so.lim.Wait(context.Background()); err != nil {
		logrus.Errorf("Failed to wait for rate limiter. Error: %v", err)
	}

	logrus.Debug("Storing apartment")
	k := datastore.NameKey("apartment", "stringID", nil)
	if _, err := so.client.Put(so.ctx, k, &ap); err != nil {
		logrus.Errorf("Failed to store apartment. Err: %v", err)
	}
}
