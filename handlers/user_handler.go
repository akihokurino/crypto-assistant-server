package handlers

import (
	"github.com/akihokurino/crypto-assistant-server/proto/go"
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
)

type userHandler struct {
	userRepository     repositories.UserRepository
	addressRepository  repositories.AddressRepository
	assetRepository    repositories.AssetRepository
	currencyRepository repositories.CurrencyRepository
}

func NewUserHandler(
	userRepository repositories.UserRepository,
	addressRepository repositories.AddressRepository,
	assetRepository repositories.AssetRepository,
	currencyRepository repositories.CurrencyRepository) pb.UserService {
	return &userHandler{
		userRepository:     userRepository,
		addressRepository:  addressRepository,
		assetRepository:    assetRepository,
		currencyRepository: currencyRepository,
	}
}

func (h *userHandler) GetAll(ctx context.Context, req *pb.Empty) (*pb.UserListResponse, error) {
	users, err := h.userRepository.GetAll(ctx)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return toUserListResponse(users), nil
}

func (h *userHandler) Get(ctx context.Context, req *pb.UserID) (*pb.UserResponse, error) {
	user, err := h.userRepository.Get(ctx, models.UserID(req.UserId))
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return toUserResponse(user), nil
}

func (h *userHandler) GetFollows(ctx context.Context, req *pb.UserID) (*pb.UserListResponse, error) {
	users, err := h.userRepository.GetFollows(ctx, models.UserID(req.UserId))
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return toUserListResponse(users), nil
}

func (h *userHandler) GetFollowers(ctx context.Context, req *pb.UserID) (*pb.UserListResponse, error) {
	users, err := h.userRepository.GetFollowers(ctx, models.UserID(req.UserId))
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return toUserListResponse(users), nil
}

func (h *userHandler) GetPortfolios(ctx context.Context, req *pb.UserID) (*pb.PortfolioListResponse, error) {
	uid := models.UserID(req.UserId)

	currencies, err := h.currencyRepository.GetAll(ctx)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	addresses, err := h.addressRepository.GetByUser(ctx, uid)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	assets, err := h.assetRepository.GetByUser(ctx, uid)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	portfolios := models.CalcPortfolios(addresses, assets, currencies, false)

	return toPortfolioListResponse(portfolios), nil
}
