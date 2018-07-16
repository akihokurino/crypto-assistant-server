package rtdb

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"github.com/akihokurino/crypto-assistant-server/utils"
	"github.com/akihokurino/crypto-assistant-server/domain/services"
	"context"
)

type PortfolioDAO struct {
	Amount   float64 `json:"amount"`
	JPYAsset float64 `json:"jpyAsset"`
}

func newPortfolioDAO(portfolio *models.Portfolio) *PortfolioDAO {
	return &PortfolioDAO{
		Amount:   portfolio.Amount,
		JPYAsset: portfolio.JPYAsset,
	}
}

func newTotalPortfolioDAO(amount float64) *PortfolioDAO {
	return &PortfolioDAO{
		Amount: amount,
	}
}

func totalPortfolioPath(userId models.UserID) string {
	return "portfolio/" + string(userId) + "/totalAmount"
}

func eachPortfolioPath(userId models.UserID, code models.CurrencyCode) string {
	return "portfolio/" + string(userId) + "/" + string(code)
}

type portfolioProvider struct {
	firebaseUtil utils.FirebaseUtil
}

func NewPortfolioProvider(firebaseUtil utils.FirebaseUtil) services.PortfolioProvider {
	return &portfolioProvider{
		firebaseUtil: firebaseUtil,
	}
}

func (p *portfolioProvider) Provide(ctx context.Context, userId models.UserID, portfolios []*models.Portfolio) error {
	client := p.firebaseUtil.InitRTDBClient(ctx)
	var lastError error
	var totalAmount float64

	for _, v := range portfolios {
		portfolioDAO := newPortfolioDAO(v)
		eachRef := client.NewRef(eachPortfolioPath(v.UserID, v.CurrencyCode))
		if err := eachRef.Set(ctx, portfolioDAO); err != nil {
			lastError = err
			break
		}
		totalAmount += v.Amount
	}

	if lastError != nil {
		return lastError
	}

	totalRef := client.NewRef(totalPortfolioPath(userId))
	if err := totalRef.Set(ctx, newTotalPortfolioDAO(totalAmount)); err != nil {
		return err
	}

	return nil
}
