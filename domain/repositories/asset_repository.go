package repositories

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"context"
)

type AssetRepository interface {
	Get(ctx context.Context, userId models.UserID, addressId models.AddressID) (*models.Asset, error)
	Put(ctx context.Context, asset *models.Asset) error
	Delete(ctx context.Context, userId models.UserID, addressId models.AddressID) error
}