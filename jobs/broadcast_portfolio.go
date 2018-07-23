package jobs

import (
	"net/http"
	"github.com/akihokurino/crypto-assistant-server/applications"
)

type broadcastPortfolio struct {
	portfolioApplication applications.PortfolioApplication
}

func NewBroadcastPortfolio(portfolioApplication applications.PortfolioApplication) JobExecutor {
	return &broadcastPortfolio{
		portfolioApplication: portfolioApplication,
	}
}

func (j *broadcastPortfolio) Exec(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := j.portfolioApplication.BroadcastAllUser(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}