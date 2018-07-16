package models

type Portfolio struct {
	UserID       UserID
	CurrencyCode CurrencyCode
	Amount       float64
	JPYAsset     float64
}

func newPortfolio(userId UserID, code CurrencyCode, amount float64, jpyAsset float64) *Portfolio {
	return &Portfolio{
		UserID:       userId,
		CurrencyCode: code,
		Amount:       amount,
		JPYAsset:     jpyAsset,
	}
}

func CalcOtherPortfolios(
	userId UserID,
	addresses []*Address,
	currencies []*Currency) []*Portfolio {

	var portfolios []*Portfolio
	for _, v1 := range currencies {
		for _, v2 := range addresses {
			if v1.Code == v2.CurrencyCode {
				portfolios = append(portfolios, newPortfolio(userId, v1.Code, 0, 0))
				break
			}
		}
	}
	return portfolios
}

func CalcMyPortfolios(
	userId UserID,
	addresses []*Address,
	assets []*Asset,
	currencies []*Currency,
	convertToJPY func(code CurrencyCode, amount float64) float64) []*Portfolio {

	amountMap := make(map[CurrencyCode]float64, len(currencies))

	for _, v := range currencies {
		amountMap[v.Code] = 0.0
	}

	for _, v1 := range addresses {
		var asset *Asset
		for _, v2 := range assets {
			if v1.Id == v2.AddressId {
				asset = v2
				break
			}
		}

		amountMap[v1.CurrencyCode] += asset.Amount
	}

	var portfolios []*Portfolio
	for code, amount := range amountMap {
		portfolios = append(portfolios, newPortfolio(userId, code, amount, convertToJPY(code, amount)))
	}
	return portfolios
}
