package store

import (
	"time"

	"github.com/wdantuma/signalk-server-go/signalk"
)

type Value struct {
	Vessel     string
	Path       string
	Value      interface{}
	Source     *signalk.Source
	Meta       *signalk.Meta
	LastChange int64
}

type Store interface {
	Put(key string, timestamp time.Time, vessel string, path string, source *signalk.Source, meta *signalk.Meta, value interface{})
	Get(key string) (*Value, bool)
	GetList(key string) []*Value
}
