package p

import (
	"context"

	"github.com/gmenih341/scrappr/scrape"
)

func Scrape(ctx context.Context, message interface{}) error {
	scrape.ScrapeHouses(ctx)

	return nil
}
