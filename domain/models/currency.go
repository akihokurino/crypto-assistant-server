package models

import "net/url"

type CurrencyCode string

type Currency struct {
	Code    CurrencyCode
	Name    string
	IconURL *url.URL
}

func NewCurrency(code CurrencyCode, name string, iconURL *url.URL) *Currency {
	return &Currency{
		Code:    code,
		Name:    name,
		IconURL: iconURL,
	}
}
