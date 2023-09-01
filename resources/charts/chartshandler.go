package charts

import (
	"encoding/json"
	"net/http"

	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type chartsHandler struct {
	chartsPath string
}

func NewChartsHandler(chartsPath string) *chartsHandler {
	return &chartsHandler{chartsPath: chartsPath}
}

func (s *chartsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	charts := make(map[string]interface{})
	path := filepath.Join(s.chartsPath, "charts")
	filepath.WalkDir(path, func(fp string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry != nil {
			info, err := entry.Info()
			if err != nil {
				return nil
			}
			name := info.Name()
			if strings.ToUpper(name) == "METADATA.JSON" {
				jsonFile, err := os.Open(fp)
				if err == nil {
					metaData := ChartMetaData{}
					bytes, err := ioutil.ReadAll(jsonFile)
					if err != nil {
						log.Print(err)
					}
					err = json.Unmarshal(bytes, &metaData)
					if err == nil {
						chart := make(map[string]interface{})
						chart["identifier"] = metaData.Id
						chart["name"] = metaData.Name
						chart["description"] = metaData.Description
						chart["format"] = metaData.Format
						chart["type"] = metaData.Type
						chart["minzoom"] = metaData.MinZoom
						chart["maxzoom"] = metaData.MaxZoom
						chart["url"] = fmt.Sprintf("/charts/%s/{z}/{x}/{y}.pbf", metaData.Id)
						charts[metaData.Id] = chart
					} else {
						log.Print(err)
					}
				}
			}
		} else {
			log.Println(fmt.Sprintf("Invalid path:%s", s.chartsPath))
		}

		return nil
	})

	chartBytes, _ := json.Marshal(charts)
	w.Write(chartBytes)
}
