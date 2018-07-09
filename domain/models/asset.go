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

func CalcAmount(assets []*Asset, prices []*CurrencyPrice) float64 {
	var amount float64
	for i, v := range assets {
		if v != nil && prices[i] != nil {
			amount += prices[i].JPY * v.Amount
		}
	}

	return amount
}
