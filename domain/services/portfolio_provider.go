package services

import (
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
)

type PortfolioProvider interface {
	Provide(ctx context.Context, userId models.UserID, portfolios []*models.Portfolio) error
} 
