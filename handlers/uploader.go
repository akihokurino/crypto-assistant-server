package handlers

import (
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"net/url"
	"google.golang.org/appengine/blobstore"
	"context"
)

type Uploader interface {
	CreateUserIconUploadURL(ctx context.Context, userId models.UserID) (*url.URL, error)
}

type uploader struct {
	userResourcePath     string
	userIconCallbackPath string
}

func NewUploader(userResourcePath, userIconCallbackPath string) Uploader {
	return &uploader{
		userResourcePath:     userResourcePath,
		userIconCallbackPath: userIconCallbackPath,
	}
}

func (u *uploader) getUserResourcePath(userId models.UserID) string {
	return u.userResourcePath + string(userId)
}

func (u *uploader) getUserIconCallbackPath(userId models.UserID) string {
	return u.userIconCallbackPath + "?uid=" + string(userId)
}

func (u *uploader) CreateUserIconUploadURL(ctx context.Context, userId models.UserID) (*url.URL, error) {
	option := blobstore.UploadURLOptions{
		MaxUploadBytes: 1024 * 1024 * 1024,
		StorageBucket:  u.getUserResourcePath(userId),
	}

	return blobstore.UploadURL(ctx, u.getUserIconCallbackPath(userId), &option)
}
