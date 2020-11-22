package realestate

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

type service struct {
	client *datastore.Client
	ctx    context.Context
}

func NewService(ctx context.Context) *service {
	projectID := os.Getenv("GCP_PROJECT")
	logrus.Infof("Running in project %s", projectID)

	dsClient, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		logrus.Fatalf("Failed to create datastore client. Err: %v", err)
	}

	return &service{dsClient, ctx}
}

func (svc *service) hasApartment(ID string) bool {
	re := &RealestateEntity{}
	key := datastore.NameKey(re.entityName(), ID, nil)
	if err := svc.client.Get(svc.ctx, key, re); err != nil {
		if err.Error() == "datastore: no such entity" {
			return false
		}
		logrus.Errorf("Failed to get realestate by ID %s. Error: %v", ID, err)
		return false
	}

	return re != nil
}

func (svc *service) hasPrice(ID string, value float32) bool {
	var rp PriceEntity
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

func (svc *service) storePrice(re RealestateEntity) {
	if !svc.hasPrice(re.ID, re.Price) {
		logrus.Infof("[DS] Storing price for %s, %f", re.ID, re.Price)

		parent := datastore.NameKey(re.entityName(), re.ID, nil)
		rp := PriceEntity{re.ID, re.Price, re.Date}
		key := datastore.NameKey(rp.entityName(), "", parent)

		if _, err := svc.client.Put(svc.ctx, key, &rp); err != nil {
			logrus.Errorf("Failed to store price. Err: %v", err)
		}
	}
}

func (svc *service) StoreRealestate(re RealestateEntity) {
	if !svc.hasApartment(re.ID) {
		logrus.Infof("[DS] Storing [%s]: %s", re.ID, re.Title)

		key := datastore.NameKey(re.entityName(), re.ID, nil)
		if _, err := svc.client.Put(svc.ctx, key, &re); err != nil {
			logrus.Errorf("Failed to store relestate. Err: %v", err)
		}
	}

	svc.storePrice(re)
}

func (svc *service) queryPricesBetweenDates(from, to time.Time) ([]PriceEntity, error) {
	query := datastore.NewQuery(PriceEntity{}.entityName()).
		Filter("createdAt >=", from).
		Filter("createdAt <=", to)

	var prices []PriceEntity
	it := svc.client.Run(svc.ctx, query)
	for {
		var price PriceEntity
		_, err := it.Next(&price)
		if err == iterator.Done {
			break
		}

		prices = append(prices, price)
	}

	return prices, nil
}

func (svc *service) GetDailyPrices() (DailyPriceResponse, error) {
	// prices, err := svc.queryPricesBetweenDates(time.Now())
	// if err != nil {
	// 	return DailyPriceResponse{}, err
	// }
	return DailyPriceResponse{}, nil
}
