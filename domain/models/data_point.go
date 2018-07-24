package models

import "time"

type DataPoint struct {
	Datetime time.Time
	USD      float64
	JPY      float64
}

func NewDataPoint(at time.Time, usd float64, jpy float64) *DataPoint {
	return &DataPoint{
		Datetime: at,
		USD:      usd,
		JPY:      jpy,
	}
}

func ConvertFromCurrencyPrice(prices []*CurrencyPrice) []*DataPoint {
	items := make([]*DataPoint, len(prices))
	for i, v := range prices {
		items[i] = NewDataPoint(v.Datetime, v.USD, v.JPY)
	}
	return items
}
