package datastore

import (
	"context"
	"google.golang.org/appengine/datastore"
)

func Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return datastore.RunInTransaction(ctx, fn, nil)
}
