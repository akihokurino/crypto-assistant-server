package services

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"context"
)

type CurrencyPriceProvider interface {
	Provide(ctx context.Context, prices []*models.CurrencyPrice) error
}