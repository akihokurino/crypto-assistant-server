package repositories

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"context"
)

type CurrencyRepository interface {
	GetAll(ctx context.Context) ([]*models.Currency, error)
	Put(ctx context.Context, currency *models.Currency) error
}