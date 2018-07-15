package jobs

import "net/http"

type JobExecutor interface {
	Exec(w http.ResponseWriter, r *http.Request)
}
