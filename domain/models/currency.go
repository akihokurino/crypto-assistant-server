package models

type CurrencyCode string

type Currency struct {
	Code CurrencyCode
	Name string
}

func NewCurrency(code CurrencyCode, name string) *Currency {
	return &Currency{
		Code: code,
		Name: name,
	}
}