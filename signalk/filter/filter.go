package filter

import (
	"github.com/wdantuma/signalk-server-go/signalk"
)

type Subscribe int

const (
	Self Subscribe = iota + 1
	All
	None
)

func ParseSubscribe(subscribe string) Subscribe {
	switch subscribe {
	case "all":
		return All
	case "none":
		return None
	default:
		return Self
	}
}

type Filter struct {
	Subscribe Subscribe
	Self      string
}

func NewFilter(self string) *Filter {
	return &Filter{Subscribe: Self, Self: self}
}

func (f *Filter) Filter(input <-chan signalk.DeltaJson) <-chan signalk.DeltaJson {
	output := make(chan signalk.DeltaJson)
	go func() {
		for delta := range input {
			include := false
			if delta.Context != nil && f.Subscribe == Self && *delta.Context == f.Self {
				include = true
			}
			if f.Subscribe == All {
				include = true
			}

			if include {
				output <- delta
			}
		}
		close(output)
	}()

	return output
}
