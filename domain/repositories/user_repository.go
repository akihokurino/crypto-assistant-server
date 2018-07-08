package repositories

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"context"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]*models.User, error)
	Get(ctx context.Context, id models.UserID) (*models.User, error)
	Put(ctx context.Context, user *models.User) error

	AlreadyFollowed(ctx context.Context, fromUserId models.UserID, toUserId models.UserID) (bool, error)
	Follow(ctx context.Context, fromUserId models.UserID, toUserId models.UserID) error
	UnFollow(ctx context.Context, fromUserId models.UserID, toUserId models.UserID) error
	GetFollows(ctx context.Context, userId models.UserID) ([]*models.User, error)
	GetFollowers(ctx context.Context, userId models.UserID) ([]*models.User, error)
}
