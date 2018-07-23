package applications

import (
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"github.com/akihokurino/crypto-assistant-server/domain/services"
)

type PortfolioApplication interface {
	BroadcastAllUser(ctx context.Context) error
	Broadcast(ctx context.Context, userId models.UserID) error
}

type portfolioApplication struct {
	userRepository          repositories.UserRepository
	currencyRepository      repositories.CurrencyRepository
	addressRepository       repositories.AddressRepository
	assetRepository         repositories.AssetRepository
	currencyPriceRepository repositories.CurrencyPriceRepository
	portfolioProvider       services.PortfolioProvider
}

func NewPortfolioApplication(
	userRepository repositories.UserRepository,
	currencyRepository repositories.CurrencyRepository,
	addressRepository repositories.AddressRepository,
	assetRepository repositories.AssetRepository,
	currencyPriceRepository repositories.CurrencyPriceRepository,
	portfolioProvider services.PortfolioProvider) PortfolioApplication {
	return &portfolioApplication{
		userRepository:          userRepository,
		currencyRepository:      currencyRepository,
		addressRepository:       addressRepository,
		assetRepository:         assetRepository,
		currencyPriceRepository: currencyPriceRepository,
		portfolioProvider:       portfolioProvider,
	}
}

func (a *portfolioApplication) BroadcastAllUser(ctx context.Context) error {
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

		portfolios := models.CalcMyPortfolios(
			v.Id,
			addresses,
			assets,
			currencies,
			func(code models.CurrencyCode, amount float64) float64 {
				price, err := a.currencyPriceRepository.GetLastByCurrency(ctx, code)
				if err != nil {
					return 0.0
				}
				return amount * price.JPY
			})

		if err := a.portfolioProvider.Provide(ctx, v.Id, portfolios); err != nil {
			lastError = err
			break
		}
	}

	return lastError
}

func (a *portfolioApplication) Broadcast(ctx context.Context, userId models.UserID) error {
	user, err := a.userRepository.Get(ctx, userId)
	if err != nil {
		return err
	}

	currencies, err := a.currencyRepository.GetAll(ctx)
	if err != nil {
		return err
	}

	addresses, err := a.addressRepository.GetByUser(ctx, user.Id)
	if err != nil {
		return err
	}

	assets, err := a.assetRepository.GetByUser(ctx, user.Id)
	if err != nil {
		return err
	}

	portfolios := models.CalcMyPortfolios(
		user.Id,
		addresses,
		assets,
		currencies,
		func(code models.CurrencyCode, amount float64) float64 {
			price, err := a.currencyPriceRepository.GetLastByCurrency(ctx, code)
			if err != nil {
				return 0.0
			}
			return amount * price.JPY
		})

	if err := a.portfolioProvider.Provide(ctx, user.Id, portfolios); err != nil {
		return err
	}

	return nil
}
