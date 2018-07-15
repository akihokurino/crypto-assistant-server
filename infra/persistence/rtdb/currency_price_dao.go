package rtdb

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"github.com/akihokurino/crypto-assistant-server/utils"
	"github.com/akihokurino/crypto-assistant-server/domain/services"
	"context"
)

type CurrencyPriceDAO struct {
	JPY float64 `json:"jpy"`
	USD float64 `json:"usd"`
}

func newCurrencyPriceDAO(price *models.CurrencyPrice) *CurrencyPriceDAO {
	return &CurrencyPriceDAO{
		JPY: price.JPY,
		USD: price.USD,
	}
}

func currencyPricePath(code models.CurrencyCode) string {
	return "currency/" + string(code)
}

type currencyPriceProvider struct {
	firebaseUtil utils.FirebaseUtil
}

func NewCurrencyPriceProvider(firebaseUtil utils.FirebaseUtil) services.CurrencyPriceProvider {
	return &currencyPriceProvider{
		firebaseUtil: firebaseUtil,
	}
}

func (s *currencyPriceProvider) Provide(ctx context.Context, prices []*models.CurrencyPrice) error {
	client := s.firebaseUtil.InitRTDBClient(ctx)
	var lastError error

	for _, v := range  prices {
		currencyPriceDAO := newCurrencyPriceDAO(v)
		ref := client.NewRef(currencyPricePath(v.CurrencyCode))
		if err := ref.Set(ctx, currencyPriceDAO); err != nil {
			lastError = err
			break
		}
	}

	return lastError
}
