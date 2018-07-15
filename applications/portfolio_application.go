package applications

import (
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"github.com/akihokurino/crypto-assistant-server/domain/services"
)

type PortfolioApplication interface {
	Broadcast(ctx context.Context) error
}

type portfolioApplication struct {
	userRepository     repositories.UserRepository
	currencyRepository repositories.CurrencyRepository
	addressRepository  repositories.AddressRepository
	assetRepository    repositories.AssetRepository
	portfolioProvider  services.PortfolioProvider
}

func NewPortfolioApplication(
	userRepository repositories.UserRepository,
	currencyRepository repositories.CurrencyRepository,
	addressRepository repositories.AddressRepository,
	assetRepository repositories.AssetRepository,
	portfolioProvider services.PortfolioProvider) PortfolioApplication {
	return &portfolioApplication{
		userRepository:     userRepository,
		currencyRepository: currencyRepository,
		addressRepository:  addressRepository,
		assetRepository:    assetRepository,
		portfolioProvider:  portfolioProvider,
	}
}

func (a *portfolioApplication) Broadcast(ctx context.Context) error {
	users, err := a.userRepository.GetAll(ctx)
	if err != nil {
		return err
	}

	currencies, err := a.currencyRepository.GetAll(ctx)
	if err != nil {
		return err
	}

	var lastError error

	for _, v := range users {
		addresses, err := a.addressRepository.GetByUser(ctx, v.Id)
		if err != nil {
			lastError = err
			break
		}

		assets, err := a.assetRepository.GetByUser(ctx, v.Id)
		if err != nil {
			lastError = err
			break
		}

		portfolios := models.CalcPortfolios(v.Id, addresses, assets, currencies, true)

		if err := a.portfolioProvider.Provide(ctx, v.Id, portfolios); err != nil {
			lastError = err
			break
		}
	}

	return lastError
}
