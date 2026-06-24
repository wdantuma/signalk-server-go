package source

import (
	"reflect"

	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/source/base"
	"go.einride.tech/can"
)

type Sourcehub struct {
	inputs  []base.DeltaSource
	output  chan signalk.DeltaJson
	started bool
}

func NewSourceHub() *Sourcehub {
	sourceHub := &Sourcehub{}
	sourceHub.output = make(chan signalk.DeltaJson)

	return sourceHub
}

func (sh *Sourcehub) AddSource(source base.DeltaSource) {
	sh.inputs = append(sh.inputs, source)
}

func (sh *Sourcehub) Sources() []string {
	sources := make([]string, 0)
	for _, s := range sh.inputs {
		sources = append(sources, s.Label())
	}
	return sources
}

func FrameValue(value interface{}) can.Frame {
	switch v := value.(type) {
	case can.Frame:
		return v
	default:
		return can.Frame{}
	}
}

func (sh *Sourcehub) Start() <-chan signalk.DeltaJson {
	return sh.startInternal()
}

func (sh *Sourcehub) startInternal() <-chan signalk.DeltaJson {
	if !sh.started {
		go func() {
			for {
				cases := make([]reflect.SelectCase, len(sh.inputs))
				for i, ch := range sh.inputs {
					cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch.Source())}
				}
				// selected, value, ok := reflect.Select(cases)
				_, value, ok := reflect.Select(cases)
				if ok {
					// input := sh.inputs[selected]
					switch v := value.Interface().(type) {
					case signalk.DeltaJson:
						sh.output <- v
					default:
						break
					}
				} else {
					break
				}
			}
			close(sh.output)
		}()
		sh.started = true
		// start inputs
		for _, input := range sh.inputs {
			input.Start()
		}
	}
	return sh.output
}
