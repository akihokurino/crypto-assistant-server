package applications

import (
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"github.com/akihokurino/crypto-assistant-server/utils"
	"errors"
	"github.com/akihokurino/crypto-assistant-server/infra/topic"
	"github.com/akihokurino/crypto-assistant-server/infra/persistence/datastore"
)

type AddressApplication interface {
	Create(ctx context.Context, userId models.UserID, currencyCode models.CurrencyCode, address string) (*models.Address, error)
	Update(ctx context.Context, userId models.UserID, addressId models.AddressID, addressText string) (*models.Address, error)
	Delete(ctx context.Context, userId models.UserID, addressId models.AddressID) error
}

type addressApplication struct {
	addressRepository     repositories.AddressRepository
	assetRepository       repositories.AssetRepository
	transactionRepository repositories.TransactionRepository
	pubsubClient          topic.PubsubClient
	idUtil                utils.IDUtil
}

func NewAddressApplication(
	addressRepository repositories.AddressRepository,
	assetRepository repositories.AssetRepository,
	transactionRepository repositories.TransactionRepository,
	pubsubClient topic.PubsubClient,
	idUtil utils.IDUtil) AddressApplication {
	return &addressApplication{
		addressRepository:     addressRepository,
		assetRepository:       assetRepository,
		transactionRepository: transactionRepository,
		pubsubClient:          pubsubClient,
		idUtil:                idUtil,
	}
}

func (a *addressApplication) Create(
	ctx context.Context,
	userId models.UserID,
	currencyCode models.CurrencyCode,
	addressText string) (*models.Address, error) {
	yes, err := a.addressRepository.ExistAddress(ctx, addressText)
	if err != nil {
		return nil, err
	}

	if yes {
		return nil, errors.New("address already created")
	}

	var address *models.Address

	if err := datastore.Transaction(ctx, func(ctx context.Context) error {
		address = models.NewAddress(models.AddressID(a.idUtil.MakeRandomKey()), userId, currencyCode, addressText)

		if err := a.addressRepository.Put(ctx, address); err != nil {
			return err
		}

		asset := models.NewAsset(userId, address.Id, 0)

		if err := a.assetRepository.Put(ctx, asset); err != nil {
			return err
		}

		transaction := models.NewTransaction(userId, address.Id, "")

		if err := a.transactionRepository.Put(ctx, transaction); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if err := a.pubsubClient.SendAddress(ctx, []*models.Address{address}); err != nil {
		return nil, err
	}

	return address, nil
}

func (a *addressApplication) Update(ctx context.Context, userId models.UserID, addressId models.AddressID, addressText string) (*models.Address, error) {
	address, err := a.addressRepository.Get(ctx, addressId)
	if err != nil {
		return nil, err
	}

	if address.UserId != userId {
		return nil, errors.New("only owner can update")
	}

	address.Value = addressText

	if err := datastore.Transaction(ctx, func(ctx context.Context) error {
		if err := a.addressRepository.Put(ctx, address); err != nil {
			return err
		}

		asset := models.NewAsset(userId, address.Id, 0)

		if err := a.assetRepository.Put(ctx, asset); err != nil {
			return err
		}

		transaction := models.NewTransaction(userId, address.Id, "")

		if err := a.transactionRepository.Put(ctx, transaction); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if err := a.pubsubClient.SendAddress(ctx, []*models.Address{address}); err != nil {
		return nil, err
	}

	return address, nil
}

func (a *addressApplication) Delete(ctx context.Context, userId models.UserID, addressId models.AddressID) error {
	address, err := a.addressRepository.Get(ctx, addressId)
	if err != nil {
		return err
	}

	if address.UserId != userId {
		return errors.New("only owner can delete")
	}

	if err := datastore.Transaction(ctx, func(ctx context.Context) error {
		if err := a.addressRepository.Delete(ctx, addressId); err != nil {
			return err
		}

		if err := a.assetRepository.Delete(ctx, userId, addressId); err != nil {
			return err
		}

		if err := a.transactionRepository.Delete(ctx, userId, addressId); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
