package api

import (
	"github.com/twitchtv/twirp"
	"net/http"
	"sort"
	"fmt"
	"strings"
	"github.com/pkg/errors"
	"context"
	logger "google.golang.org/appengine/log"
	"github.com/akihokurino/crypto-assistant-server/utils"
)

func newLoggingServerHooks() *twirp.ServerHooks {
	contextUtil := utils.NewContextUtil()

	headerJoiner := func(src http.Header) string {
		var resSrc []string
		var sortedKeySet []string
		for k := range src {
			sortedKeySet = append(sortedKeySet, k)
		}
		sort.Strings(sortedKeySet)
		for _, k := range sortedKeySet {
			for _, v := range src[k] {
				resSrc = append(resSrc, fmt.Sprintf("    %s : %s", k, v))
			}
		}
		return strings.Join(resSrc, "\n")
	}
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			req, ok := contextUtil.HttpRequest(ctx)
			if !ok {
				return ctx, errors.WithStack(errors.New("request info has gone."))
			}
			var logStrSrc []string
			logStrSrc = append(logStrSrc, "\nurl: " + req.URL.String())
			logStrSrc = append(logStrSrc, "host: " + req.Host)
			logStrSrc = append(logStrSrc, "method: " + req.Method)
			logStrSrc = append(logStrSrc, "header-info:\n" + headerJoiner(req.Header))

			logger.Infof(ctx, strings.Join(logStrSrc, "\n"))
			return ctx, nil
		},
		ResponseSent: func(ctx context.Context) {
			w, ok := contextUtil.HttpResponseWriter(ctx)
			if ok {
				var logStrSrc []string
				logStrSrc = append(logStrSrc, "\nheader-info:\n"+headerJoiner(w.Header()))
				logger.Infof(ctx, strings.Join(logStrSrc, "\n"))
			}
		},
	}
}