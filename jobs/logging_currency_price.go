package jobs

import (
	"net/http"
	"github.com/akihokurino/crypto-assistant-server/applications"
)

type LoggingCurrencyPrice interface {
	Exec(w http.ResponseWriter, r *http.Request)
}

type loggingCurrencyPrice struct {
	currencyPriceApplication applications.CurrencyPriceApplication
}

func NewLoggingCurrencyPrice(currencyPriceApplication applications.CurrencyPriceApplication) LoggingCurrencyPrice {
	return &loggingCurrencyPrice{
		currencyPriceApplication: currencyPriceApplication,
	}
}

func (j *loggingCurrencyPrice) Exec(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := j.currencyPriceApplication.CreateEachTime(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
