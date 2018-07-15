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

	userRepository := datastore.NewUserRepository(dateUtil)
	currencyRepository := datastore.NewCurrencyRepository()
	currencyPriceRepository := datastore.NewCurrencyPriceRepository()
	addressRepository := datastore.NewAddressRepository()
	assetRepository := datastore.NewAssetRepository()

	currencyPriceProvider := rtdb.NewCurrencyPriceProvider(firebaseUtil)
	portfolioProvider := rtdb.NewPortfolioProvider(firebaseUtil)

	currencyApplication := applications.NewCurrencyApplication(currencyRepository)
	currencyPriceApplication := applications.NewCurrencyPriceApplication(
		httpClient,
		currencyRepository,
		currencyPriceRepository,
		currencyPriceProvider,
		idUtil,
		dateUtil)
	portfolioApplication := applications.NewPortfolioApplication(
		userRepository,
		currencyRepository,
		addressRepository,
		assetRepository,
		portfolioProvider,
	)

	mux.HandleFunc(
		"/job/register_currencies",
		appEngine(jobs.NewRegisterCurrencies(currencyApplication).Exec))

	mux.HandleFunc(
		"/job/update_market",
		appEngine(jobs.NewUpdateMarket(
			addressRepository,
			currencyPriceApplication,
			pubsubClient).Exec))

	mux.HandleFunc(
		"/job/broadcast_portfolio",
		appEngine(jobs.NewBroadcastPortfolio(portfolioApplication).Exec))
}
