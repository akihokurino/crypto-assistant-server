package handlers

import (
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"github.com/akihokurino/crypto-assistant-server/proto/go"
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
)

type chartHandler struct {
	currencyPriceRepository repositories.CurrencyPriceRepository
}

func NewChartHandler(currencyPriceRepository repositories.CurrencyPriceRepository) pb.ChartService {
	return &chartHandler{
		currencyPriceRepository: currencyPriceRepository,
	}
}

func (h *chartHandler) GetLast24Hour(ctx context.Context, req *pb.CurrencyCode) (*pb.ChartResponse, error) {
	prices, err := h.currencyPriceRepository.GetLast24HourByCurrency(ctx, models.CurrencyCode(req.CurrencyCode))
	if err != nil {
		return nil, handleError(ctx, err)
	}

	points := models.ConvertFromCurrencyPrice(prices)

	return toChartResponse(points), nil
}
