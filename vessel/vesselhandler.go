package vessel

import (
	"net/http"

	"github.com/wdantuma/signalk-server-go/signalkserver/state"
)

type vesselHandler struct {
	state state.ServerState
}

func NewVesselHandler(s state.ServerState) *vesselHandler {
	return &vesselHandler{state: s}
}

func (s *vesselHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
