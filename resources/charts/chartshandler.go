package charts

import (
	"encoding/json"
	"net/http"
)

type chartsHandler struct {
}

func NewChartsHandler() *chartsHandler {
	return &chartsHandler{}
}

func (s *chartsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	c1 := make(map[string]interface{})
	c1["identifier"] = "test"
	c1["name"] = "S-57 test"
	c1["description"] = "20230721_U7Inland_Waddenzee_week 29_NL"
	c1["format"] = "pbf"
	c1["type"] = "S-57"
	c1["minZoom"] = 14
	c1["maxZoom"] = 14
	c1["url"] = "http://localhost:3000/charts/test/{x}/{y}/{z}"

	charts := make(map[string]interface{})

	charts["c1"] = c1

	chartBytes, _ := json.Marshal(charts)
	w.Write(chartBytes)

}
