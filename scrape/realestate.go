package scrape

import "time"

type realestate struct {
	ID               string    `datastore:"id"`
	Title            string    `datastore:"title"`
	Price            float32   `datastore:"-"`
	Area             float32   `datastore:"area"`
	Image            string    `datastore:"image,noindex"`
	ShortDescription string    `datastore:"shortDescription,noindex"`
	LongDescription  string    `datastore:"longDescription,noindex"`
	Date             time.Time `datastore:"date"`
}

func (re realestate) entityName() string {
	return "realestate"
}

func (re realestate) toRow() []interface{} {
	return []interface{}{
		re.ID,
		re.Title,
		re.Price,
		re.Area,
		re.Image,
		re.ShortDescription,
		re.LongDescription,
		re.Date.Format("02.01.2006"),
	}
}

type realestatePrice struct {
	ID    string  `datastore:"id"`
	Price float32 `datastore:"price"`
	Date  time.Time
}

func (rp realestatePrice) entityName() string {
	return "realestatePrice"
}
