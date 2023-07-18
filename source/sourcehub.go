package source

import (
	"reflect"

	"go.einride.tech/can"
)

type Sourcehub struct {
	inputs  []CanSource
	output  chan SourceFrame
	started bool
}

func (sh *Sourcehub) AddSource(source CanSource) {
	sh.inputs = append(sh.inputs, source)
	sh.startInternal()
}

func (sh *Sourcehub) Sources() []string {
	sources := make([]string, 0)
	for _, s := range sh.inputs {
		sources = append(sources, s.Label())
	}
	return sources
}

func NewSourceHub() *Sourcehub {
	sourceHub := &Sourcehub{}
	sourceHub.output = make(chan SourceFrame)
	return sourceHub
}

func FrameValue(value interface{}) can.Frame {
	switch v := value.(type) {
	case can.Frame:
		return v
	default:
		return can.Frame{}
	}
}

func (sh *Sourcehub) Start() <-chan SourceFrame {
	return sh.output
}

func (sh *Sourcehub) startInternal() <-chan SourceFrame {
	if !sh.started {
		go func() {
			for {
				cases := make([]reflect.SelectCase, len(sh.inputs))
				for i, ch := range sh.inputs {
					cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch.Source())}
				}
				selected, value, ok := reflect.Select(cases)
				if ok {
					input := sh.inputs[selected]
					sourceFrame := SourceFrame{Frame: FrameValue(value.Interface()), Label: input.Label()}
					sh.output <- sourceFrame

				} else {
					break
				}
			}
			close(sh.output)
		}()
		sh.started = true
	}
	return sh.output
}
