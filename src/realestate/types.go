package realestate

import "time"

type PriceTuple struct {
	Price      PriceEntity
	Realestate RealestateEntity
}
type DailyPriceResponse map[time.Time]PriceTuple
