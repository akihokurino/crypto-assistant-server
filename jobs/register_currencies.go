package jobs

import (
	"net/http"
	"github.com/akihokurino/crypto-assistant-server/applications"
)

type RegisterCurrencies interface {
	Exec(w http.ResponseWriter, r *http.Request)
}

type registerCurrencies struct {
	currencyApplication applications.CurrencyApplication
}

func NewRegisterCurrencies(currencyApplication applications.CurrencyApplication) RegisterCurrencies {
	return &registerCurrencies{
		currencyApplication: currencyApplication,
	}
}

func (j *registerCurrencies) Exec(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := j.currencyApplication.Create(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}