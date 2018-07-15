package models

type Portfolio struct {
	UserID       UserID
	CurrencyCode CurrencyCode
	Amount       float64
}

func newPortfolio(userId UserID, currency *Currency, amount float64) *Portfolio {
	return &Portfolio{
		UserID:       userId,
		CurrencyCode: currency.Code,
		Amount:       amount,
	}
}

func CalcPortfolios(
	userId UserID,
	addresses []*Address,
	assets []*Asset,
	currencies []*Currency,
	showAmount bool) []*Portfolio {
	portfolios := make([]*Portfolio, len(currencies))
	for i, v1 := range currencies {
		var address *Address
		for _, v2 := range addresses {
			if v1.Code == v2.CurrencyCode {
				address = v2
				break
			}
		}

		if address == nil {
			portfolios[i] = newPortfolio(userId, v1, 0)
			continue
		}

		var asset *Asset
		for _, v2 := range assets {
			if address.Id == v2.AddressId {
				asset = v2
				break
			}
		}

		if asset == nil {
			portfolios[i] = newPortfolio(userId, v1, 0)
			continue
		}

		if showAmount {
			portfolios[i] = newPortfolio(userId, v1, asset.Amount)
		} else {
			portfolios[i] = newPortfolio(userId, v1, 0)
		}
	}
	return portfolios
}
