package charts

import (
	"net/http"
)

type chartsHandler struct {
}

func NewChartsHandler() *chartsHandler {
	return &chartsHandler{}
}

func (s *chartsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotFound)

}
