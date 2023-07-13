package vessel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/signalkserver/state"
	"github.com/wdantuma/signalk-server-go/store"
)

type vesselHandler struct {
	state state.ServerState
}

func NewVesselHandler(s state.ServerState) *vesselHandler {
	return &vesselHandler{state: s}
}

func MapValue(value interface{}) map[string]interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		return v
	default:
		return nil
	}
}

func GetResultObject(r map[string]interface{}, parts []string, v *store.Value) map[string]interface{} {
	if r == nil {
		r = make(map[string]interface{})
	}
	if len(parts) == 1 {
		r[parts[0]] = v.Value
	} else {
		r2 := MapValue(r[parts[0]])
		if r2 == nil {
			r2 = make(map[string]interface{})
		}
		r[parts[0]] = GetResultObject(r2, parts[1:], v)
	}
	return r
}

func (s *vesselHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	metaRequest := false
	key := ""
	parts := strings.Split(r.RequestURI, "/")[5:]
	if len(parts) > 0 && parts[0] != "" {
		if r.RequestURI[len(r.RequestURI)-1:] == "/" {
			parts = parts[:len(parts)-1]
		}
		vessel := parts[0]
		parts = parts[1:]
		if vessel == "self" {
			vessel = s.state.GetSelf()
		}
		if len(parts) > 1 && parts[len(parts)-1] == "meta" {
			metaRequest = true
			parts = parts[:len(parts)-1]
		}
		path := strings.Join(parts, ".")
		key = fmt.Sprintf("%s/%s", vessel, path)
	}

	if metaRequest {
		result := make(map[string]interface{})

		resultBytes, _ := json.Marshal(result)

		w.Write(resultBytes)
	} else {
		values := s.state.GetStore().GetList(key)

		if len(values) > 0 {
			result := signalk.SignalkJson{Vessels: make(signalk.SignalkJsonVessels)}
			result.Self = s.state.GetSelf()
			var resultVessel map[string]interface{}
			lastVessel := ""
			for _, v := range values {
				if v.Vessel != lastVessel {
					if resultVessel != nil {
						result.Vessels[lastVessel] = resultVessel
					}
					lastVessel = v.Vessel
					resultVessel = nil
				}
				parts := strings.Split(v.Path, ".")
				resultVessel = GetResultObject(resultVessel, parts, v)
			}
			result.Vessels[lastVessel] = resultVessel

			resultBytes, _ := json.Marshal(result)

			w.Write(resultBytes)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
