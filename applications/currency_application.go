package applications

import (
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
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

	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("ETH"), "Ethereum"))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("ETC"), "Ethereum Classic"))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("EOS"), "EOS"))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("TRX"), "TRON"))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("BNB"), "Binance Coin"))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("OMG"), "OmiseGO"))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("VEN"), "VeChain"))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("ICX"), "ICON"))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("ZIL"), "Zilliqa"))
	currencies = append(currencies, models.NewCurrency(models.CurrencyCode("ZRX"), "Ox"))

	for _, currency := range currencies {
		_ = a.currencyRepository.Put(ctx, currency)
	}

	return nil
}
