package models

import "time"

type CurrencyPriceID string

type CurrencyPrice struct {
	Id           CurrencyPriceID
	CurrencyCode CurrencyCode
	USD          float64
	JPY          float64
	Datetime     time.Time
}

func NewCurrencyPrice(id CurrencyPriceID, code CurrencyCode, usd float64, jpy float64, now time.Time) *CurrencyPrice {
	return &CurrencyPrice{
		Id:           id,
		CurrencyCode: code,
		USD:          usd,
		JPY:          jpy,
		Datetime:     now,
	}
}
