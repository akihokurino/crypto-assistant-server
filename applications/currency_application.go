package applications

import (
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"net/url"
)

type CurrencyApplication interface {
	Create(ctx context.Context) error
}

type currencyApplication struct {
	currencyRepository repositories.CurrencyRepository
}

func NewCurrencyApplication(currencyRepository repositories.CurrencyRepository) CurrencyApplication {
	return &currencyApplication{
		currencyRepository: currencyRepository,
	}
}

func (a *currencyApplication) Create(ctx context.Context) error {
	var currencies []*models.Currency

	iconURL, _ := url.Parse("https://storage.cloud.google.com/crypto-assistant-dev.appspot.com/currencies/ethereum.png")
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("ETH"), "Ethereum", iconURL))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("ETC"), "Ethereum Classic", iconURL))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("EOS"), "EOS", iconURL))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("TRX"), "TRON", iconURL))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("BNB"), "Binance Coin", iconURL))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("OMG"), "OmiseGO", iconURL))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("VEN"), "VeChain", iconURL))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("ICX"), "ICON", iconURL))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("ZIL"), "Zilliqa", iconURL))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("ZRX"), "Ox", iconURL))

	for _, currency := range currencies {
		_ = a.currencyRepository.Put(ctx, currency)
	}

	return nil
}
