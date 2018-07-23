package api

import (
	"net/http"
	"github.com/akihokurino/crypto-assistant-server/proto/go"
	"github.com/akihokurino/crypto-assistant-server/handlers"
	"github.com/akihokurino/crypto-assistant-server/infra/persistence/datastore"
	"github.com/akihokurino/crypto-assistant-server/applications"
	"github.com/akihokurino/crypto-assistant-server/utils"
	"github.com/akihokurino/crypto-assistant-server/infra/topic"
	"github.com/akihokurino/crypto-assistant-server/infra/persistence/rtdb"
)

const uploadUserIconCallbackPath = "/users/icon/upload_callback"

func init() {
	mux := http.DefaultServeMux

	contextUtil := utils.NewContextUtil()
	idUtil := utils.NewIDUtil()
	dateUtil := utils.NewDateUtil()
	firebaseUtil := utils.NewFirebaseUtil()

	pubsubClient := topic.NewPubsubClient()

	currencyRepository := datastore.NewCurrencyRepository()
	currencyPriceRepository := datastore.NewCurrencyPriceRepository()
	userRepository := datastore.NewUserRepository(dateUtil)
	addressRepository := datastore.NewAddressRepository()
	assetRepository := datastore.NewAssetRepository()
	transactionRepository := datastore.NewTransactionRepository()

	portfolioProvider := rtdb.NewPortfolioProvider(firebaseUtil)

	userApplication := applications.NewUserApplication(userRepository)
	uploadApplication := applications.NewUploadApplication(userRepository)
	addressApplication := applications.NewAddressApplication(
		addressRepository,
		assetRepository,
		transactionRepository,
		pubsubClient,
		idUtil)
	portfolioApplication := applications.NewPortfolioApplication(
		userRepository,
		currencyRepository,
		addressRepository,
		assetRepository,
		currencyPriceRepository,
		portfolioProvider,
	)

	uploader := handlers.NewUploader(
		"crypto-assistant-dev.appspot.com/users/",
		uploadUserIconCallbackPath)

	currencyService := pb.NewCurrencyServiceServer(handlers.NewCurrencyHandler(
		currencyRepository), newLoggingServerHooks())
	mux.Handle(
		pb.CurrencyServicePathPrefix,
		cros(appEngine(currencyService)))

	currencyPriceService := pb.NewCurrencyPriceServiceServer(handlers.NewCurrencyPriceHandler(
		currencyRepository,
		currencyPriceRepository), newLoggingServerHooks())
	mux.Handle(
		pb.CurrencyPriceServicePathPrefix,
		cros(appEngine(currencyPriceService)))

	meService := pb.NewMeServiceServer(handlers.NewMeHandler(
		userRepository,
		addressRepository,
		assetRepository,
		currencyPriceRepository,
		currencyRepository,
		userApplication,
		uploader,
		contextUtil), newLoggingServerHooks())
	mux.Handle(pb.MeServicePathPrefix, cros(appEngine(auth(meService))))

	userService := pb.NewUserServiceServer(handlers.NewUserHandler(
		userRepository,
		addressRepository,
		currencyRepository), newLoggingServerHooks())
	mux.Handle(pb.UserServicePathPrefix, cros(appEngine(auth(userService))))

	addressService := pb.NewAddressServiceServer(handlers.NewAddressHandler(
		addressApplication,
		portfolioApplication,
		contextUtil), newLoggingServerHooks())
	mux.Handle(pb.AddressServicePathPrefix, cros(appEngine(auth(addressService))))

	mux.HandleFunc(uploadUserIconCallbackPath, handlers.NewUploadCallbackHandler(uploadApplication).CallbackUserIcon)
}
