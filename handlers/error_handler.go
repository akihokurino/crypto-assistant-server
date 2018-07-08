package handlers

import (
	"github.com/twitchtv/twirp"
	"context"
	logger "google.golang.org/appengine/log"
)

func handleError(ctx context.Context, err error) error {
	return toInternalServerError(ctx, err)
}

func toInternalServerError(ctx context.Context, err error) error {
	logger.Infof(ctx, err.Error())
	return twirp.InternalErrorWith(err)
}
