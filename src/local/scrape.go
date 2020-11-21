package main

import (
	"context"
	"encoding/json"

	p "github.com/gmenih341/scrappr/src/functions"
	"github.com/gmenih341/scrappr/src/scrape"
)

func main() {
	bytes, _ := json.Marshal(scrape.FunctionMessage{"https://www.nepremicnine.net/oglasi-prodaja/podravska/maribor/stanovanje/", "apartment"})
	p.Scrape(context.Background(), p.PubSubMessage{bytes})
}
