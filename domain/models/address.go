package models

type AddressID string

type Address struct {
	Id           AddressID
	UserId       UserID
	CurrencyCode CurrencyCode
	Value        string
}

func NewAddress(id AddressID, userId UserID, currencyCode CurrencyCode, value string) *Address {
	return &Address{
		Id: id,
		UserId: userId,
		CurrencyCode: currencyCode,
		Value: value,
	}
}