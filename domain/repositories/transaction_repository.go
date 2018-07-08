package repositories

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"context"
)

type TransactionRepository interface {
	Put(ctx context.Context, transaction *models.Transaction) error
	Delete(ctx context.Context, userId models.UserID, addressId models.AddressID) error
}
