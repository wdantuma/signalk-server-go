package charts

import (
	"time"
)

type ChartMetaData struct {
	Id           string    `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	Created      time.Time `json:"created,omitempty"`
	ChartUpdated time.Time `json:"chartupdated,omitempty"`
	ChartType    string    `json:"charttype,omitempty"`
	ChartFormat  string    `json:"chartformat,omitempty"`
}
