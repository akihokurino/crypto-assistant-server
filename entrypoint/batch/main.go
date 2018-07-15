package batch

import (
	"net/http"
	"github.com/akihokurino/crypto-assistant-server/jobs"
	"github.com/akihokurino/crypto-assistant-server/infra/api"
	"github.com/akihokurino/crypto-assistant-server/infra/persistence/datastore"
	"github.com/akihokurino/crypto-assistant-server/applications"
	"github.com/akihokurino/crypto-assistant-server/utils"
	"github.com/akihokurino/crypto-assistant-server/infra/topic"
	"github.com/akihokurino/crypto-assistant-server/infra/persistence/rtdb"
)

func init() {
	mux := http.DefaultServeMux

	idUtil := utils.NewIDUtil()
	dateUtil := utils.NewDateUtil()
	firebaseUtil := utils.NewFirebaseUtil()

	httpClient := api.NewHttpClient()
	pubsubClient := topic.NewPubsubClient()

	currencyRepository := datastore.NewCurrencyRepository()
	currencyPriceRepository := datastore.NewCurrencyPriceRepository()
	addressRepository := datastore.NewAddressRepository()

	currencyPriceProvider := rtdb.NewCurrencyPriceProvider(firebaseUtil)

	currencyApplication := applications.NewCurrencyApplication(currencyRepository)

	currencyPriceApplication := applications.NewCurrencyPriceApplication(
		httpClient,
		currencyRepository,
		currencyPriceRepository,
		currencyPriceProvider,
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
