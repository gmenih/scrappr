package scrape

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"google.golang.org/api/iterator"
)

type datastoreSvc struct {
	*rate.Limiter
	client *datastore.Client
	ctx    context.Context
}

func newStore(ctx context.Context) *datastoreSvc {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	logrus.Infof("Running in project %s", projectID)

	dsClient, err := datastore.NewClient(ctx, projectID)
	// n, m = n requests every m seconds
	limiter := rate.NewLimiter(50, 1)

	if err != nil {
		logrus.Fatalf("Failed to create sheets client. Err: %v", err)
	}

	return &datastoreSvc{limiter, dsClient, ctx}
}

func (svc *datastoreSvc) wait() {
	if !svc.Allow() {
		logrus.Warning("Waiting for Google API limit")
	}

	if err := svc.Wait(svc.ctx); err != nil {
		logrus.Errorf("Failed to wait for rate limiter. Error: %v", err)
	}
}

func (svc *datastoreSvc) findByID(ID string) (*realestate, error) {
	svc.wait()

	re := &realestate{}
	key := datastore.NameKey(re.entityName(), ID, nil)
	if err := svc.client.Get(svc.ctx, key, re); err != nil {
		if err.Error() == "datastore: no such entity" {
			return nil, nil
		}
		logrus.Errorf("Failed to get realestate by ID %s. Error: %v", ID, err)
		return nil, errors.New("Failed to get realestate by ID")
	}

	return re, nil
}

func (svc *datastoreSvc) entityHasPrice(ID string, value float32) bool {
	var rp price
	query := datastore.NewQuery(rp.entityName()).
		Filter("id =", ID).
		Filter("price =", value).
		KeysOnly()

	it := svc.client.Run(svc.ctx, query)
	_, err := it.Next(&rp)

	if err == nil || err == iterator.Done {
		return false
	} else if err != nil {
		logrus.Errorf("Failed to check price existence. Error: %v", err)
	}

	return true
}

func (svc *datastoreSvc) storePrice(re realestate) {
	if !svc.entityHasPrice(re.ID, re.Price) {
		logrus.Infof("[DS] Storing price for %s, %f", re.ID, re.Price)
		parent := datastore.NameKey(re.entityName(), re.ID, nil)
		rp := price{re.ID, re.Price, re.Date}
		key := datastore.NameKey(rp.entityName(), "", parent)

		if _, err := svc.client.Put(svc.ctx, key, &rp); err != nil {
			logrus.Errorf("Failed to store price. Err: %v", err)
		}
	}
}

func (svc *datastoreSvc) storeRealestate(re realestate) {
	svc.wait()

	logrus.Debug("[DS] Storing [%s]: %s", re.ID, re.Title)

	matched, err := svc.findByID(re.ID)
	if err != nil {
		logrus.Errorf("Failed to find an existing relestate")
	}

	if matched == nil {
		logrus.Infof("Storing relestate %s", re.ID)
		key := datastore.NameKey(re.entityName(), re.ID, nil)
		if _, err := svc.client.Put(svc.ctx, key, &re); err != nil {
			logrus.Errorf("Failed to store relestate. Err: %v", err)
		}
	}

	svc.storePrice(re)
}
