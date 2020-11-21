package p

import (
	"context"
	"encoding/json"

	"github.com/gmenih341/scrappr/src/scrape"
	"github.com/sirupsen/logrus"
)

func Graph(ctx context.Context, message PubSubMessage) error {
	body := scrape.FunctionMessage{}
	if err := json.Unmarshal(message.Data, &body); err != nil {
		logrus.New().Panicf("Failed to parse message!")
	}

	// graph.Graph(ctx, body)

	return nil
}
