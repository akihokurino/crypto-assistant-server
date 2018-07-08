package handlers

import (
	"github.com/akihokurino/crypto-assistant-server/proto/go"
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
)

type userHandler struct {
	userRepository repositories.UserRepository
}

func NewUserHandler(userRepository repositories.UserRepository) pb.UserService {
	return &userHandler{
		userRepository: userRepository,
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
