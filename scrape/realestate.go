package scrape

import "time"

type RealestateType string

type realestate struct {
	ID               string    `datastore:"id"`
	Title            string    `datastore:"title"`
	Price            float32   `datastore:"-"`
	Area             float32   `datastore:"area"`
	Image            string    `datastore:"image,noindex"`
	ShortDescription string    `datastore:"shortDescription,noindex"`
	LongDescription  string    `datastore:"longDescription,noindex"`
	Type             string    `datastore:"type"`
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

type price struct {
	ID        string    `datastore:"id"`
	Price     float32   `datastore:"price"`
	CreatedAt time.Time `datastore:"createdAt"`
}

func (rp price) entityName() string {
	return "price"
}
