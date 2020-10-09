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
	id               int32
	title            string
	priceCents       float32
	area             float32
	image            string
	shortDescription string
	longDescription  string
	date             time.Time
}

func (a house) toRow() []interface{} {
	return []interface{}{
		a.id,
		a.title,
		a.priceCents,
		a.area,
		a.image,
		a.shortDescription,
		a.longDescription,
		a.date.Format("02.01.2006"),
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
	sheet := newOp(ctx)

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
			title:            title,
			shortDescription: shortDescription,
			longDescription:  longDescription,
			image:            thumbnailURL,
			area:             parseArea(shortDescription),
			date:             time.Now(),
		}

		fmt.Sscanf(price, "%f", &a.priceCents)
		fmt.Sscanf(idString, "%d", &a.id)
		if strings.Contains(priceText, "m2") {
			a.priceCents = a.priceCents * a.area
		}

		logrus.Infof("Storing house %d", a.id)
		sheet.storeApartment(a)
	})

	if err := pageScraper.Visit(siteURL); err != nil {
		logrus.Errorf("Failed to visit %s! Error: %v", siteURL, err)
	}

	pageScraper.Wait()
	apartmentScraper.Wait()
	logrus.Infof("Found %d apts", count)
}
