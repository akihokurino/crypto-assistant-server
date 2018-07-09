package batch

import (
	"net/http"
	"github.com/akihokurino/crypto-assistant-server/jobs"
	"github.com/akihokurino/crypto-assistant-server/infra/api"
	"github.com/akihokurino/crypto-assistant-server/infra/persistence/datastore"
	"github.com/akihokurino/crypto-assistant-server/applications"
	"github.com/akihokurino/crypto-assistant-server/utils"
	"github.com/akihokurino/crypto-assistant-server/infra/topic"
)

func init() {
	mux := http.DefaultServeMux

	idUtil := utils.NewIDUtil()
	dateUtil := utils.NewDateUtil()

	httpClient := api.NewHttpClient()
	pubsubClient := topic.NewPubsubClient()

	currencyRepository := datastore.NewCurrencyRepository()
	currencyPriceRepository := datastore.NewCurrencyPriceRepository()
	addressRepository := datastore.NewAddressRepository()

	currencyApplication := applications.NewCurrencyApplication(currencyRepository)

	currencyPriceApplication := applications.NewCurrencyPriceApplication(
		httpClient,
		currencyRepository,
		currencyPriceRepository,
		idUtil,
		dateUtil)

	mux.HandleFunc(
		"/job/logging_currency_price",
		appEngine(jobs.NewLoggingCurrencyPrice(currencyPriceApplication).Exec))

	mux.HandleFunc(
		"/job/register_currencies",
		appEngine(jobs.NewRegisterCurrencies(currencyApplication).Exec))

	mux.HandleFunc(
		"/job/update_asset",
		appEngine(jobs.NewUpdateAsset(addressRepository, pubsubClient).Exec))
}
