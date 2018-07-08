package handlers

import (
	"github.com/akihokurino/crypto-assistant-server/proto/go"
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
)

type currencyHandler struct {
	currencyRepository repositories.CurrencyRepository
}

func NewCurrencyHandler(currencyRepository repositories.CurrencyRepository) pb.CurrencyService {
	return &currencyHandler{
		currencyRepository: currencyRepository,
	}
}

func (h *currencyHandler) GetAll(ctx context.Context, req *pb.Empty) (*pb.CurrencyListResponse, error) {
	currencies, err := h.currencyRepository.GetAll(ctx)

	if err != nil {
		return nil, handleError(ctx, err)
	}

	return toCurrencyListResponse(currencies), nil
}