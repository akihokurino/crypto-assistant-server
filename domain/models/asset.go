package models

type Asset struct {
	UserId    UserID
	AddressId AddressID
	Amount    float64
}

func NewAsset(userId UserID, addressId AddressID, amount float64) *Asset {
	return &Asset{
		UserId: userId,
		AddressId: addressId,
		Amount: amount,
	}
}
