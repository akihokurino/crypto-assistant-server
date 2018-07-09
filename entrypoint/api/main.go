package api

import (
	"net/http"
	"github.com/akihokurino/crypto-assistant-server/proto/go"
	"github.com/akihokurino/crypto-assistant-server/handlers"
	"github.com/akihokurino/crypto-assistant-server/infra/persistence/datastore"
	"github.com/akihokurino/crypto-assistant-server/applications"
	"github.com/akihokurino/crypto-assistant-server/utils"
	"github.com/akihokurino/crypto-assistant-server/infra/topic"
)

const uploadUserIconCallbackPath = "/users/icon/upload_callback"

func init() {
	mux := http.DefaultServeMux

	contextUtil := utils.NewContextUtil()
	idUtil := utils.NewIDUtil()
	dateUtil := utils.NewDateUtil()

	pubsubClient := topic.NewPubsubClient()

	currencyRepository := datastore.NewCurrencyRepository()
	currencyPriceRepository := datastore.NewCurrencyPriceRepository()
	userRepository := datastore.NewUserRepository(dateUtil)
	addressRepository := datastore.NewAddressRepository()
	assetRepository := datastore.NewAssetRepository()
	transactionRepository := datastore.NewTransactionRepository()

	userApplication := applications.NewUserApplication(userRepository)
	uploadApplication := applications.NewUploadApplication(userRepository)
	addressApplication := applications.NewAddressApplication(
		addressRepository,
		assetRepository,
		transactionRepository,
		pubsubClient,
		idUtil)

	uploader := handlers.NewUploader(
		"crypto-assistant-dev.appspot.com/users/",
		uploadUserIconCallbackPath)

	currencyService := pb.NewCurrencyServiceServer(handlers.NewCurrencyHandler(
		currencyRepository), newLoggingServerHooks())
	mux.Handle(
		pb.CurrencyServicePathPrefix,
		appEngine(auth(currencyService)))

	currencyPriceService := pb.NewCurrencyPriceServiceServer(handlers.NewCurrencyPriceHandler(
		currencyRepository,
		currencyPriceRepository), newLoggingServerHooks())
	mux.Handle(
		pb.CurrencyPriceServicePathPrefix,
		appEngine(auth(currencyPriceService)))

	meService := pb.NewMeServiceServer(handlers.NewMeHandler(
		userRepository,
		addressRepository,
		assetRepository,
		currencyPriceRepository,
		currencyRepository,
		userApplication,
		uploader,
		contextUtil), newLoggingServerHooks())
	mux.Handle(pb.MeServicePathPrefix, appEngine(auth(meService)))

	userService := pb.NewUserServiceServer(handlers.NewUserHandler(
		userRepository,
		addressRepository,
		assetRepository,
		currencyRepository), newLoggingServerHooks())
	mux.Handle(pb.UserServicePathPrefix, appEngine(auth(userService)))

	addressService := pb.NewAddressServiceServer(handlers.NewAddressHandler(
		addressApplication,
		contextUtil), newLoggingServerHooks())
	mux.Handle(pb.AddressServicePathPrefix, appEngine(auth(addressService)))

	mux.HandleFunc(uploadUserIconCallbackPath, handlers.NewUploadCallbackHandler(uploadApplication).CallbackUserIcon)
}
