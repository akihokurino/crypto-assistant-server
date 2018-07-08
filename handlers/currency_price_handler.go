package handlers

import (
	"github.com/akihokurino/crypto-assistant-server/proto/go"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
)

type currencyPriceHandler struct {
	currencyRepository      repositories.CurrencyRepository
	currencyPriceRepository repositories.CurrencyPriceRepository
}

func NewCurrencyPriceHandler(
	currencyRepository repositories.CurrencyRepository,
	currencyPriceRepository repositories.CurrencyPriceRepository) pb.CurrencyPriceService {
	return &currencyPriceHandler{
		currencyRepository:      currencyRepository,
		currencyPriceRepository: currencyPriceRepository,
	}
}

func (h *currencyPriceHandler) GetLast(ctx context.Context, req *pb.Empty) (*pb.CurrencyPriceListResponse, error) {
	currencies, err := h.currencyRepository.GetAll(ctx)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	prices := make([]*models.CurrencyPrice, len(currencies))
	for i, v := range currencies {
		price, err := h.currencyPriceRepository.GetLastByCurrency(ctx, v.Code)
		if err != nil {
			continue
		}
		prices[i] = price
	}

	return toCurrencyPriceListResponse(prices), nil
}

func (h *currencyPriceHandler) GetByCurrency(ctx context.Context, req *pb.CurrencyCode) (*pb.CurrencyPriceListResponse, error) {
	prices, err := h.currencyPriceRepository.GetByCurrency(ctx, models.CurrencyCode(req.CurrencyCode))
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return toCurrencyPriceListResponse(prices), nil
}
