package applications

import (
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"github.com/pkg/errors"
)

type UserApplication interface {
	Create(ctx context.Context, userId models.UserID, username string) (*models.User, error)
	Update(ctx context.Context, userId models.UserID, username string) (*models.User, error)
	Follow(ctx context.Context, fromUserId models.UserID, toUserId models.UserID) error
	UnFollow(ctx context.Context, fromUserId models.UserID, toUserId models.UserID) error
}

type userApplication struct {
	userRepository repositories.UserRepository
}

func NewUserApplication(userRepository repositories.UserRepository) UserApplication {
	return &userApplication{
		userRepository: userRepository,
	}
}

func (a *userApplication) Create(ctx context.Context, userId models.UserID, username string) (*models.User, error) {
	user := models.NewUser(userId, username)

	if err := a.userRepository.Put(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (a *userApplication) Update(ctx context.Context, userId models.UserID, username string) (*models.User, error) {
	user := models.NewUser(userId, username)

	if err := a.userRepository.Put(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (a *userApplication) Follow(ctx context.Context, fromUserId models.UserID, toUserId models.UserID) error {
	yes, err := a.userRepository.AlreadyFollowed(ctx, fromUserId, toUserId)
	if err != nil {
		return err
	}

	if yes {
		return errors.New("already followed")
	}

	if err := a.userRepository.Follow(ctx, fromUserId, toUserId); err != nil {
		return err
	}

	return nil
}

func (a *userApplication) UnFollow(ctx context.Context, fromUserId models.UserID, toUserId models.UserID) error {
	if err := a.userRepository.UnFollow(ctx, fromUserId, toUserId); err != nil {
		return err
	}

	return nil
}