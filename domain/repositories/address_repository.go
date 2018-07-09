package repositories

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"context"
)

type AddressRepository interface {
	GetAll(ctx context.Context) ([]*models.Address, error)
	GetByUser(ctx context.Context, userId models.UserID) ([]*models.Address, error)
	Get(ctx context.Context, addressId models.AddressID) (*models.Address, error)
	ExistAddress(ctx context.Context, addressText string) (bool, error)
	Put(ctx context.Context, address *models.Address) error
	Delete(ctx context.Context, addressId models.AddressID) error
}