package jobs

import (
	"net/http"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"github.com/akihokurino/crypto-assistant-server/infra/topic"
)

type UpdateAsset interface {
	Exec(w http.ResponseWriter, r *http.Request)
}

type updateAsset struct {
	addressRepository repositories.AddressRepository
	pubsubClient topic.PubsubClient
}

func NewUpdateAsset(addressRepository repositories.AddressRepository, pubsubClient topic.PubsubClient) UpdateAsset {
	return &updateAsset{
		addressRepository: addressRepository,
		pubsubClient: pubsubClient,
	}
}

func (j *updateAsset) Exec(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
