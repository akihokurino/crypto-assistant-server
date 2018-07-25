package applications

import (
	"context"
	"github.com/akihokurino/crypto-assistant-server/infra/api"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"github.com/akihokurino/crypto-assistant-server/utils"
	"strings"
	"os"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"github.com/akihokurino/crypto-assistant-server/domain/services"
)

type CurrencyPriceApplication interface {
	CreateEachTime(ctx context.Context) error
}

type currencyPriceApplication struct {
	httpClient              api.HttpClient
	currencyRepository      repositories.CurrencyRepository
	currencyPriceRepository repositories.CurrencyPriceRepository
	currencyPriceProvider   services.CurrencyPriceProvider
	idUtil                  utils.IDUtil
	dateUtil                utils.DateUtil
}

func NewCurrencyPriceApplication(
	httpClient api.HttpClient,
	currencyRepository repositories.CurrencyRepository,
	currencyPriceRepository repositories.CurrencyPriceRepository,
	currencyPriceProvider services.CurrencyPriceProvider,
	idUtil utils.IDUtil,
	dateUtil utils.DateUtil) CurrencyPriceApplication {
	return &currencyPriceApplication{
		httpClient:              httpClient,
		currencyRepository:      currencyRepository,
		currencyPriceRepository: currencyPriceRepository,
		currencyPriceProvider:   currencyPriceProvider,
		idUtil:                  idUtil,
		dateUtil:                dateUtil,
	}
}

func (a *currencyPriceApplication) CreateEachTime(ctx context.Context) error {
	now := a.dateUtil.CurrentTime()

	currencies, err := a.currencyRepository.GetAll(ctx)
	if err != nil {
		return err
	}

	var codes []string
	for _, v := range currencies {
		codes = append(codes, string(v.Code))
	}

	params := map[string]string{
		"fsyms": strings.Join(codes, ","),
		"tsyms": "USD,JPY",
	}

	res, err := a.httpClient.Get(ctx, os.Getenv("CRYPTO_COMPARE_ENDPOINT"), params)

	if err != nil {
		return err
	}

	prices := make([]*models.CurrencyPrice, len(currencies))

	for i, v := range currencies {
		priceMap := res[string(v.Code)].(map[string]interface{})
		usd := priceMap["USD"].(float64)
		jpy := priceMap["JPY"].(float64)

		currencyPrice := models.NewCurrencyPrice(
			models.CurrencyPriceID(a.idUtil.MakeRandomKey()),
			v.Code,
			usd,
			jpy,
			now)

		prices[i] = currencyPrice

		_ = a.currencyPriceRepository.Put(ctx, currencyPrice)
	}

	if err := a.currencyPriceProvider.Provide(ctx, prices); err != nil {
		return err
	}

	return nil
}
