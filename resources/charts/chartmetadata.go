package charts

import (
	"time"
)

type ChartMetaData struct {
	Id          string    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Created     time.Time `json:"created,omitempty"`
	Updated     time.Time `json:"updated,omitempty"`
	Type        string    `json:"type,omitempty"`
	Format      string    `json:"format,omitempty"`
	MinZoom     int       `json:"minzoom,omitempty"`
	MaxZoom     int       `json:"maxzoom,omitempty"`
	Bounds      []float32 `json:"bounds,omitempty"`
	Scale       int       `json:"scale,omitempty"`
}
