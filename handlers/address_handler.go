package handlers

import (
	"github.com/akihokurino/crypto-assistant-server/proto/go"
	"context"
	"github.com/akihokurino/crypto-assistant-server/applications"
	"github.com/akihokurino/crypto-assistant-server/utils"
	"errors"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
)

type addressHandler struct {
	addressApplication   applications.AddressApplication
	portfolioApplication applications.PortfolioApplication
	contextUtil          utils.ContextUtil
}

func NewAddressHandler(
	addressApplication applications.AddressApplication,
	portfolioApplication applications.PortfolioApplication,
	contextUtil utils.ContextUtil) pb.AddressService {
	return &addressHandler{
		addressApplication:   addressApplication,
		portfolioApplication: portfolioApplication,
		contextUtil:          contextUtil,
	}
}

func (h *addressHandler) Create(ctx context.Context, req *pb.CreateAddressRequest) (*pb.AddressResponse, error) {
	uid, ok := h.contextUtil.AuthUID(ctx)
	if !ok {
		return nil, handleError(ctx, errors.New("failed resolve dependency"))
	}

	address, err := h.addressApplication.Create(ctx, uid, models.CurrencyCode(req.CurrencyCode), req.Value)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	if err := h.portfolioApplication.Broadcast(ctx, uid); err != nil {
		return nil, handleError(ctx, err)
	}

	return toAddressResponse(address), nil
}

func (h *addressHandler) Update(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.AddressResponse, error) {
	uid, ok := h.contextUtil.AuthUID(ctx)
	if !ok {
		return nil, handleError(ctx, errors.New("failed resolve dependency"))
	}

	address, err := h.addressApplication.Update(ctx, uid, models.AddressID(req.AddressId), req.Value)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	if err := h.portfolioApplication.Broadcast(ctx, uid); err != nil {
		return nil, handleError(ctx, err)
	}

	return toAddressResponse(address), nil
}

func (h *addressHandler) Delete(ctx context.Context, req *pb.AddressID) (*pb.Empty, error) {
	uid, ok := h.contextUtil.AuthUID(ctx)
	if !ok {
		return nil, handleError(ctx, errors.New("failed resolve dependency"))
	}

	if err := h.addressApplication.Delete(ctx, uid, models.AddressID(req.AddressId)); err != nil {
		return nil, handleError(ctx, err)
	}

	if err := h.portfolioApplication.Broadcast(ctx, uid); err != nil {
		return nil, handleError(ctx, err)
	}

	return toEmptyResponse(), nil
}
