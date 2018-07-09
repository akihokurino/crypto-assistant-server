package handlers

import (
	"github.com/akihokurino/crypto-assistant-server/proto/go"
	"context"
	"github.com/akihokurino/crypto-assistant-server/utils"
	"errors"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"github.com/akihokurino/crypto-assistant-server/applications"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
)

type meHandler struct {
	userRepository          repositories.UserRepository
	addressRepository       repositories.AddressRepository
	assetRepository         repositories.AssetRepository
	currencyPriceRepository repositories.CurrencyPriceRepository
	currencyRepository      repositories.CurrencyRepository
	userApplication         applications.UserApplication
	uploader                Uploader
	contextUtil             utils.ContextUtil
}

func NewMeHandler(
	userRepository repositories.UserRepository,
	addressRepository repositories.AddressRepository,
	assetRepository repositories.AssetRepository,
	currencyPriceRepository repositories.CurrencyPriceRepository,
	currencyRepository repositories.CurrencyRepository,
	userApplication applications.UserApplication,
	uploader Uploader,
	contextUtil utils.ContextUtil) pb.MeService {
	return &meHandler{
		userRepository:          userRepository,
		addressRepository:       addressRepository,
		assetRepository:         assetRepository,
		currencyPriceRepository: currencyPriceRepository,
		currencyRepository:      currencyRepository,
		userApplication:         userApplication,
		uploader:                uploader,
		contextUtil:             contextUtil,
	}
}

func (h *meHandler) Get(ctx context.Context, req *pb.Empty) (*pb.UserResponse, error) {
	uid, ok := h.contextUtil.AuthUID(ctx)
	if !ok {
		return nil, handleError(ctx, errors.New("failed resolve dependency"))
	}

	user, err := h.userRepository.Get(ctx, uid)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return toUserResponse(user), nil
}

func (h *meHandler) Create(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	uid, ok := h.contextUtil.AuthUID(ctx)
	if !ok {
		return nil, handleError(ctx, errors.New("failed resolve dependency"))
	}

	user, err := h.userApplication.Create(ctx, uid, req.Username)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return toUserResponse(user), nil
}

func (h *meHandler) Update(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	uid, ok := h.contextUtil.AuthUID(ctx)
	if !ok {
		return nil, handleError(ctx, errors.New("failed resolve dependency"))
	}

	user, err := h.userApplication.Update(ctx, uid, req.Username)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return toUserResponse(user), nil
}

func (h *meHandler) CreateUploadIconURL(ctx context.Context, req *pb.Empty) (*pb.UploadURL, error) {
	uid, ok := h.contextUtil.AuthUID(ctx)
	if !ok {
		return nil, handleError(ctx, errors.New("failed resolve dependency"))
	}

	uploadUrl, err := h.uploader.CreateUserIconUploadURL(ctx, uid)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return toUploadURLResponse(uploadUrl), nil
}

func (h *meHandler) GetAddresses(ctx context.Context, req *pb.Empty) (*pb.AddressListResponse, error) {
	uid, ok := h.contextUtil.AuthUID(ctx)
	if !ok {
		return nil, handleError(ctx, errors.New("failed resolve dependency"))
	}

	addresses, err := h.addressRepository.GetByUser(ctx, uid)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return toAddressListResponse(addresses), nil
}

func (h *meHandler) Follow(ctx context.Context, req *pb.UserID) (*pb.Empty, error) {
	uid, ok := h.contextUtil.AuthUID(ctx)
	if !ok {
		return nil, handleError(ctx, errors.New("failed resolve dependency"))
	}

	if err := h.userApplication.Follow(ctx, uid, models.UserID(req.UserId)); err != nil {
		return nil, handleError(ctx, err)
	}

	return toEmptyResponse(), nil
}

func (h *meHandler) UnFollow(ctx context.Context, req *pb.UserID) (*pb.Empty, error) {
	uid, ok := h.contextUtil.AuthUID(ctx)
	if !ok {
		return nil, handleError(ctx, errors.New("failed resolve dependency"))
	}

	if err := h.userApplication.UnFollow(ctx, uid, models.UserID(req.UserId)); err != nil {
		return nil, handleError(ctx, err)
	}

	return toEmptyResponse(), nil
}

func (h *meHandler) GetFollows(ctx context.Context, req *pb.Empty) (*pb.UserListResponse, error) {
	uid, ok := h.contextUtil.AuthUID(ctx)
	if !ok {
		return nil, handleError(ctx, errors.New("failed resolve dependency"))
	}

	users, err := h.userRepository.GetFollows(ctx, uid)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return toUserListResponse(users), nil
}

func (h *meHandler) GetFollowers(ctx context.Context, req *pb.Empty) (*pb.UserListResponse, error) {
	uid, ok := h.contextUtil.AuthUID(ctx)
	if !ok {
		return nil, handleError(ctx, errors.New("failed resolve dependency"))
	}

	users, err := h.userRepository.GetFollowers(ctx, uid)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	return toUserListResponse(users), nil
}

func (h *meHandler) GetAsset(ctx context.Context, req *pb.Empty) (*pb.AssetResponse, error) {
	uid, ok := h.contextUtil.AuthUID(ctx)
	if !ok {
		return nil, handleError(ctx, errors.New("failed resolve dependency"))
	}

	addresses, err := h.addressRepository.GetByUser(ctx, uid)
	if err != nil {
		return nil, handleError(ctx, err)
	}

	var amount float64
	for _, v := range addresses {
		asset, err := h.assetRepository.Get(ctx, uid, v.Id)
		if err != nil {
			continue
		}

		currentPrice, err := h.currencyPriceRepository.GetLastByCurrency(ctx, v.CurrencyCode)
		if err != nil {
			continue
		}

		amount += currentPrice.JPY * asset.Amount
	}

	return toAssetResponse(amount), nil
}

func (h *meHandler) GetPortfolios(ctx context.Context, req *pb.Empty) (*pb.PortfolioListResponse, error) {
	uid, ok := h.contextUtil.AuthUID(ctx)
	if !ok {
		return nil, handleError(ctx, errors.New("failed resolve dependency"))
	}

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

	portfolios := models.CalcPortfolios(addresses, assets, currencies, true)

	return toPortfolioListResponse(portfolios), nil
}
