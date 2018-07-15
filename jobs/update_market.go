package jobs

import (
	"net/http"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"github.com/akihokurino/crypto-assistant-server/infra/topic"
	"github.com/akihokurino/crypto-assistant-server/applications"
)

type updateMarket struct {
	addressRepository        repositories.AddressRepository
	currencyPriceApplication applications.CurrencyPriceApplication
	pubsubClient             topic.PubsubClient
}

func NewUpdateMarket(
	addressRepository repositories.AddressRepository,
	currencyPriceApplication applications.CurrencyPriceApplication,
	pubsubClient topic.PubsubClient) JobExecutor {
	return &updateMarket{
		addressRepository:        addressRepository,
		currencyPriceApplication: currencyPriceApplication,
		pubsubClient:             pubsubClient,
	}
}

func (j *updateMarket) Exec(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := j.currencyPriceApplication.CreateEachTime(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	addresses, err := j.addressRepository.GetAll(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := j.pubsubClient.SendAddress(ctx, addresses); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
