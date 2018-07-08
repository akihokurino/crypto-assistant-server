package utils

import (
	"context"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"net/http"
)

const requestStoreKey = "__request_store_key__"
const responseStoreKey = "__response_store_key__"
const authUidStoreKey = "__auth_uid_store_key__"

type ContextUtil interface {
	WithAuthUID(ctx context.Context, uid models.UserID) (context.Context, error)
	WithHTTPRequest(ctx context.Context, h http.Request) (context.Context, error)
	WithHTTPResponseWriter(ctx context.Context, w http.ResponseWriter) (context.Context, error)
	AuthUID(ctx context.Context) (models.UserID, bool)
	HttpRequest(ctx context.Context) (http.Request, bool)
	HttpResponseWriter(ctx context.Context) (http.ResponseWriter, bool)
}

type contextUtil struct {

}

func NewContextUtil() ContextUtil {
	return &contextUtil{}
}

func (u *contextUtil) WithAuthUID(ctx context.Context, uid models.UserID) (context.Context, error) {
	return context.WithValue(ctx, authUidStoreKey, uid), nil
}

func (u *contextUtil) WithHTTPRequest(ctx context.Context, h http.Request) (context.Context, error) {
	return context.WithValue(ctx, requestStoreKey, h), nil
}

func (u *contextUtil) WithHTTPResponseWriter(ctx context.Context, w http.ResponseWriter) (context.Context, error) {
	return context.WithValue(ctx, responseStoreKey, w), nil
}

func (u *contextUtil) AuthUID(ctx context.Context) (models.UserID, bool) {
	uid, ok := ctx.Value(authUidStoreKey).(models.UserID)
	return uid, ok
}

func (u *contextUtil) HttpRequest(ctx context.Context) (http.Request, bool) {
	h, ok := ctx.Value(requestStoreKey).(http.Request)
	return h, ok
}

func (u *contextUtil) HttpResponseWriter(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value(responseStoreKey).(http.ResponseWriter)
	return w, ok
}
