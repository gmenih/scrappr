package scrape

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
)

const siteURL = "https://www.nepremicnine.net/oglasi-prodaja/podravska/maribor/hisa/"
const baseURL = "https://www.nepremicnine.net"

type house struct {
	ID               int32     `datastore:"id"`
	Title            string    `datastore:"title"`
	PriceCents       float32   `datastore:"priceCents"`
	Area             float32   `datastore:"area"`
	Image            string    `datastore:"image,noindex"`
	ShortDescription string    `datastore:"shortDescription,noindex"`
	LongDescription  string    `datastore:"longDescription,noindex"`
	Date             time.Time `datastore:"date"`
}

func (a house) toRow() []interface{} {
	return []interface{}{
		a.ID,
		a.Title,
		a.PriceCents,
		a.Area,
		a.Image,
		a.ShortDescription,
		a.LongDescription,
		a.Date.Format("02.01.2006"),
	}
}

func parseArea(description string) float32 {
	r, err := regexp.Compile(`(\d+(,\d+)?) m2`)
	if err != nil {
		return 0
	}
	var area float32
	stringArea := strings.ReplaceAll(r.FindString(description), ",", ".")
	fmt.Sscanf(stringArea, "%f m2", &area)

	return area
}

func visitUrl(c *colly.Collector, apartmentURL string) {
	if err := c.Visit(baseURL + apartmentURL); err != nil {
		logrus.Errorf("Failed to visit %s, Err: %+v", apartmentURL, err)
	}
}

func ScrapeHouses(ctx context.Context) {
	logrus.Infof("Starting to scrape")
	count := 0

	pageScraper := colly.NewCollector(
		colly.UserAgent("AmigaBot (rstate v5.2)"),
	)

	pageScraper.Limit(&colly.LimitRule{
		Delay:       time.Second * 100,
		Parallelism: 1,
	})

	apartmentScraper := colly.NewCollector(
		colly.UserAgent("AmigaBot (rstate v5.2)"),
	)
	apartmentScraper.Limit(&colly.LimitRule{
		Delay:       time.Second * 100,
		Parallelism: 1,
	})
	store := newStore(ctx)

	pageScraper.OnHTML(".teksti_container[data-href]", func(el *colly.HTMLElement) {
		apartmentURL := el.Attr("data-href")

		if strings.Index(apartmentURL, "/oglasi-prodaja") != 0 {
			logrus.Warningf("Skipping bot trap link %s", apartmentURL)
			return
		}

		logrus.Infof("Matched house. Visiting %s", apartmentURL)
		if err := apartmentScraper.Visit(el.Request.AbsoluteURL(apartmentURL)); err != nil {
			logrus.Errorf("Fucking error! %v", err)
		}
	})

	pageScraper.OnHTML(".headbar > #pagination > ul > li.paging_next > a.next", func(el *colly.HTMLElement) {
		pageURL := el.Attr("href")
		logrus.Debugf("Matched page. Visiting %s", pageURL)

		el.Request.Visit(el.Request.AbsoluteURL(pageURL))
	})

	apartmentScraper.OnHTML("div[itemprop=mainEntity][id]", func(el *colly.HTMLElement) {
		idRxp := regexp.MustCompile(`(\d+)\/$`)
		groups := idRxp.FindAllStringSubmatch(el.Request.URL.Path, -1)
		if len(groups) != 1 || groups[0][1] == "" {
			return
		}

		count++

		idString := groups[0][1]
		title := el.ChildAttr("meta[itemprop=name]", "content")
		price := el.ChildAttr("meta[itemprop=price]", "content")
		priceText := el.ChildText(".cena > span")
		shortDescription := el.ChildText("div.kratek")
		thumbnailURL := el.ChildAttr("meta[itemprop=image]", "content")
		longDescription := el.ChildText("div[itemprop=disambiguatingDescription]")
		logrus.Infof("House[%s]: %s", idString, title)

		a := house{
			Title:            title,
			ShortDescription: shortDescription,
			LongDescription:  longDescription,
			Image:            thumbnailURL,
			Area:             parseArea(shortDescription),
			Date:             time.Now(),
		}

		fmt.Sscanf(price, "%f", &a.PriceCents)
		fmt.Sscanf(idString, "%d", &a.ID)
		if strings.Contains(priceText, "m2") {
			a.PriceCents = a.PriceCents * a.Area
		}

		logrus.Infof("Storing house %d", a.ID)
		store.storeApartment(a)
	})

	if err := pageScraper.Visit(siteURL); err != nil {
		logrus.Errorf("Failed to visit %s! Error: %v", siteURL, err)
	}

	pageScraper.Wait()
	apartmentScraper.Wait()
	logrus.Infof("Found %d apts", count)
}
