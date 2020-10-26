package scrape

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type datastoreSvc struct {
	*rate.Limiter
	client *datastore.Client
}

func newStore(ctx context.Context) *datastoreSvc {
	dsClient, err := datastore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	// n, m = n requests every m seconds
	limiter := rate.NewLimiter(50, 1)

	if err != nil {
		logrus.Fatalf("Failed to create sheets client. Err: %v", err)
	}

	return &datastoreSvc{limiter, dsClient}
}

func (svc *datastoreSvc) wait() {
	if !svc.Allow() {
		logrus.Warning("Waiting for Google API limit")
	}

	if err := svc.Wait(context.Background()); err != nil {
		logrus.Errorf("Failed to wait for rate limiter. Error: %v", err)
	}
}

func (svc *datastoreSvc) findByID(ID string) (*realestate, error) {
	svc.wait()

	re := &realestate{}
	key := datastore.NameKey(re.entityName(), ID, nil)
	if err := svc.client.Get(context.Background(), key, re); err != nil {
		if err.Error() == "datastore: no such entity" {
			return nil, nil
		}
		logrus.Errorf("Failed to get realestate by ID %s. Error: %v", ID, err)
		return nil, errors.New("Failed to get realestate by ID")
	}

	return re, nil
}

func (svc *datastoreSvc) storePrice(re realestate) {
	svc.wait()

	logrus.Debugf("[DS] Storing price for %s, %d", re.ID, re.Price)

	parent := datastore.NameKey(re.entityName(), re.ID, nil)
	rp := realestatePrice{re.ID, re.Price, re.Date}
	key := datastore.NameKey(rp.entityName(), "", parent)
	if _, err := svc.client.Put(context.Background(), key, &rp); err != nil {
		logrus.Errorf("Failed to store apartment. Err: %v", err)
	}
}

func (svc *datastoreSvc) storeRealestate(re realestate) {
	svc.wait()

	logrus.Debug("[DS] Storing [%s]: %s", re.ID, re.Title)

	matched, err := svc.findByID(re.ID)
	if err != nil {
		logrus.Errorf("Failed to find an existing apartment")
	}

	if matched == nil {
		key := datastore.NameKey(re.entityName(), re.ID, nil)
		if _, err := svc.client.Put(context.Background(), key, &re); err != nil {
			logrus.Errorf("Failed to store apartment. Err: %v", err)
		}
	}

	svc.storePrice(re)
}
