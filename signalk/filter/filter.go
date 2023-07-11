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

type Filter struct {
	Subscribe Subscribe
}

func NewFilter() *Filter {
	return &Filter{Subscribe: Self}
}

func (f *Filter) Filter(input <-chan signalk.DeltaJson) <-chan signalk.DeltaJson {
	output := make(chan signalk.DeltaJson)
	go func() {
		for delta := range input {
			include := true

			// for _, update := range delta.Updates {
			// 	if uint(*update.Source.Pgn) == 130306 {
			// 		include = true
			// 	}
			// }
			//newDelta := DeltaJson{}
			//delta.Context
			if include {
				output <- delta
			}
		}
		close(output)
	}()

	return output
}
