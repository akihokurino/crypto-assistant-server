package applications

import (
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"github.com/akihokurino/crypto-assistant-server/domain/repositories"
	"net/url"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/image"
	"errors"
)

type UploadApplication interface {
	UploadUserIcon(ctx context.Context, userId models.UserID, blobs map[string][]*blobstore.BlobInfo) error
}

type uploadApplication struct {
	userRepository repositories.UserRepository
}

func NewUploadApplication(userRepository repositories.UserRepository) UploadApplication {
	return &uploadApplication{
		userRepository: userRepository,
	}
}

func (a *uploadApplication) UploadUserIcon(
	ctx context.Context,
	userId models.UserID,
	blobs map[string][]*blobstore.BlobInfo) error {
	servingURL, err := getServingURL(ctx, blobs)
	if err != nil {
		return err
	}

	user, err := a.userRepository.Get(ctx, userId)
	if err != nil {
		return err
	}

	user.IconURL = servingURL

	if err := a.userRepository.Put(ctx, user); err != nil {
		return err
	}

	return nil
}

func getServingURL(ctx context.Context, blobs map[string][]*blobstore.BlobInfo) (*url.URL, error) {
	file := blobs["file"]

	if len(file) == 0 {
		return nil, errors.New("")
	}

	option := image.ServingURLOptions{Secure: false, Crop: false}
	return image.ServingURL(ctx, file[0].BlobKey, &option)
}
