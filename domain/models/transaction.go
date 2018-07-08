package models

type Transaction struct {
	UserId    UserID
	AddressId AddressID
	Text      string
}

func NewTransaction(userId UserID, addressId AddressID, text string) *Transaction {
	return &Transaction{
		UserId: userId,
		AddressId: addressId,
		Text: text,
	}
}
