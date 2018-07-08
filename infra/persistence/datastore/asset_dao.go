package datastore

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"context"
	"github.com/mjibson/goon"
	"github.com/pkg/errors"
	"google.golang.org/appengine/datastore"
)

const kindAsset = "Asset"

type AssetDAO struct {
	_kind     string `goon:"kind,Asset"`
	Id        string `datastore:"-" goon:"id"`
	UserId    string
	AddressId string
	Amount    float64
}

func onlyIdAssetDAO(userId models.UserID, addressId models.AddressID) *AssetDAO {
	return &AssetDAO{Id: string(userId) + "-" + string(addressId)}
}

func (d *AssetDAO) toModel() *models.Asset {
	return &models.Asset{
		UserId:    models.UserID(d.UserId),
		AddressId: models.AddressID(d.AddressId),
		Amount:    d.Amount,
	}
}

func fromAssetToDAO(from *models.Asset) *AssetDAO {
	return &AssetDAO{
		Id:        string(from.UserId) + "-" + string(from.AddressId),
		UserId:    string(from.UserId),
		AddressId: string(from.AddressId),
		Amount:    from.Amount,
	}
}

type assetRepository struct {
}

func NewAssetRepository() repositories.AssetRepository {
	return &assetRepository{}
}

func (r *assetRepository) Get(ctx context.Context, userId models.UserID, addressId models.AddressID) (*models.Asset, error) {
	g := goon.FromContext(ctx)

	assetDAO := onlyIdAssetDAO(userId, addressId)

	if err := g.Get(assetDAO); err != nil {
		return nil, errors.WithStack(err)
	}

	return assetDAO.toModel(), nil
}

func (r *assetRepository) Put(ctx context.Context, asset *models.Asset) error {
	g := goon.FromContext(ctx)

	if _, err := g.Put(fromAssetToDAO(asset)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *assetRepository) Delete(ctx context.Context, userId models.UserID, addressId models.AddressID) error {
	g := goon.FromContext(ctx)

	assetDAO := onlyIdAssetDAO(userId, addressId)

	key := datastore.NewKey(g.Context, kindAsset, assetDAO.Id, 0, nil)

	if err := g.Delete(key); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
