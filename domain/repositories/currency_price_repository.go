package repositories

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"context"
)

type CurrencyPriceRepository interface {
	GetByCurrency(ctx context.Context, code models.CurrencyCode) ([]*models.CurrencyPrice, error)
	GetLastByCurrency(ctx context.Context, code models.CurrencyCode) (*models.CurrencyPrice, error)
	GetLast24HourByCurrency(ctx context.Context, code models.CurrencyCode) ([]*models.CurrencyPrice, error)
	Put(ctx context.Context, currencyPrice *models.CurrencyPrice) error
}