package store

import (
	"github.com/wdantuma/signalk-server-go/signalk"
)

type Value struct {
	Vessel     string
	Path       string
	Value      interface{}
	Source     *signalk.Source
	LastChange int64
}

type Store interface {
	Put(key string, timestamp int64, vessel string, path string, source *signalk.Source, value interface{})
	Get(key string) (*Value, bool)
	GetList(key string) []*Value
}
