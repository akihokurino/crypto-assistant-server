package batch

import (
	"net/http"
	"google.golang.org/appengine"
)

func appEngine(base http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.WithContext(r.Context(), r)
		base(w, r.WithContext(ctx))
	}
}