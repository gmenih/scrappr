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

	pageColly := colly.NewCollector(
		colly.UserAgent("AmigaBot (rstate v5.2)"),
	)

	pageColly.Limit(&colly.LimitRule{
		Delay:       time.Second * 100,
		Parallelism: 1,
	})

	realestateColly := colly.NewCollector(
		// Random user agent to be a bit less obvious that we're scraping 🙈
		colly.UserAgent("AmigaBot (rstate v5.2)"),
	)
	realestateColly.Limit(&colly.LimitRule{
		Delay:       time.Second * 100,
		Parallelism: 1,
	})
	store := newStore(ctx)

	pageColly.OnHTML(".teksti_container[data-href]", func(el *colly.HTMLElement) {
		apartmentURL := el.Attr("data-href")

		if strings.Index(apartmentURL, "/oglasi-prodaja") != 0 {
			logrus.Warningf("Skipping bot trap link %s", apartmentURL)
			return
		}

		logrus.Infof("Matched house. Visiting %s", apartmentURL)
		if err := realestateColly.Visit(el.Request.AbsoluteURL(apartmentURL)); err != nil {
			logrus.Errorf("Fucking error! %v", err)
		}
	})

	pageColly.OnHTML(".headbar > #pagination > ul > li.paging_next > a.next", func(el *colly.HTMLElement) {
		pageURL := el.Attr("href")
		logrus.Debugf("Matched page. Visiting %s", pageURL)

		el.Request.Visit(el.Request.AbsoluteURL(pageURL))
	})

	realestateColly.OnHTML("div[itemprop=mainEntity][id]", func(el *colly.HTMLElement) {
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

		r := realestate{
			ID:               idString,
			Title:            title,
			ShortDescription: shortDescription,
			LongDescription:  longDescription,
			Image:            thumbnailURL,
			Area:             parseArea(shortDescription),
			Date:             time.Now(),
		}

		fmt.Sscanf(price, "%f", &r.Price)
		if strings.Contains(priceText, "m2") {
			r.Price = r.Price * r.Area
		}

		logrus.Infof("Storing realestate %s", r.ID)
		store.storeRealestate(r)
	})

	if err := pageColly.Visit(siteURL); err != nil {
		logrus.Errorf("Failed to visit %s! Error: %v", siteURL, err)
	}

	pageColly.Wait()
	realestateColly.Wait()
	logrus.Infof("Found %d apts", count)
}
