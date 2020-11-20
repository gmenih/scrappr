package p

import (
	"context"
	"encoding/json"

	"github.com/gmenih341/scrappr/scrape"
	"github.com/sirupsen/logrus"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func Scrape(ctx context.Context, message PubSubMessage) error {
	body := scrape.FunctionMessage{}
	if err := json.Unmarshal(message.Data, &body); err != nil {
		logrus.New().Panicf("Failed to parse message!")
	}
	scrape.ScrapeRealestate(ctx, body)

	return nil
}
