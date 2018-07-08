package handlers

import (
	"github.com/akihokurino/crypto-assistant-server/proto/go"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
)

type currencyPriceHandler struct {
	currencyPriceRepository repositories.CurrencyPriceRepository
}

func NewCurrencyPriceHandler(currencyPriceRepository repositories.CurrencyPriceRepository) pb.CurrencyPriceService {
	return &currencyPriceHandler{
		currencyPriceRepository: currencyPriceRepository,
	}
}

func (h *currencyPriceHandler) GetByCurrency(ctx context.Context, req *pb.CurrencyCode) (*pb.CurrencyPriceListResponse, error) {
	prices, err := h.currencyPriceRepository.GetByCurrency(ctx, models.CurrencyCode(req.CurrencyCode))

	if err != nil {
		return nil, handleError(ctx, err)
	}

	return toCurrencyPriceListResponse(prices), nil
}
