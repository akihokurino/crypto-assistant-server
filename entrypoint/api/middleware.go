package api

import (
	"net/http"
	"google.golang.org/appengine"
	"github.com/akihokurino/crypto-assistant-server/utils"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
)

func appEngine(base http.Handler) http.Handler {
	contextUtil := utils.NewContextUtil()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.WithContext(r.Context(), r)
		ctx, e := contextUtil.WithHTTPRequest(ctx, *r)
		ctx, e = contextUtil.WithHTTPResponseWriter(ctx, w)
		if e != nil {
			panic(e)
		}
		base.ServeHTTP(w, r.WithContext(ctx))
	})
}

func auth(base http.Handler) http.Handler {
	contextUtil := utils.NewContextUtil()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		//ローカルで起動した場合はアクセストークンではなくuidをつけてリクエストする
		//X-Debug-User-Id: l7W16MD993NuV9IGJYpRJxkVE543
		if appengine.IsDevAppServer() {
			uid := r.Header.Get("X-Debug-User-Id")
			newContext, _ := contextUtil.WithAuthUID(ctx, models.UserID(uid))
			base.ServeHTTP(w, r.WithContext(newContext))
		} else {
			client := utils.NewFirebaseUtil().InitAuthClient(ctx)

			token := r.Header.Get("Authorization")
			if len(token) <= 7 {
				w.WriteHeader(403)
				return
			}

			decoded, err := client.VerifyIDToken(token[7:])
			if err != nil {
				w.WriteHeader(403)
				return
			}

			newContext, _ := contextUtil.WithAuthUID(ctx, models.UserID(decoded.UID))
			base.ServeHTTP(w, r.WithContext(newContext))
		}
	})
}

func cros(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		base.ServeHTTP(w, r)
	})
}
