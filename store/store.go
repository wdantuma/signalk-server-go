package store

import (
	"github.com/wdantuma/signalk-server-go/signalk"
)

type Value struct {
	value      interface{}
	source     *signalk.Source
	lastChange int64
}

type Store interface {
	Put(key string, timestamp int64, source *signalk.Source, value interface{})
	Get(key string) (*Value, bool)
}
