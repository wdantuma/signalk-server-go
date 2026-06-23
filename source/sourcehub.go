package source

import (
	"reflect"

	// nmea0183converter "github.com/wdantuma/signalk-server-go/converter/nmea0183"
	nmea2000converter "github.com/wdantuma/signalk-server-go/converter/nmea2000"
	"github.com/wdantuma/signalk-server-go/signalk"
	nmea2000source "github.com/wdantuma/signalk-server-go/source/nmea2000"
	"go.einride.tech/can"
)

type Sourcehub struct {
	inputs       []DeltaSource
	output       chan signalk.DeltaJson
	started      bool
	n2kConverter nmea2000converter.Nmea2000ToSignalk
	// nmeaConverter nmea0183converter.Nmea0183ToSignalk

}

func NewSourceHub(n2kConverter nmea2000converter.Nmea2000ToSignalk) *Sourcehub {
	sourceHub := &Sourcehub{}
	sourceHub.output = make(chan signalk.DeltaJson)
	sourceHub.n2kConverter = n2kConverter

	return sourceHub
}

func (sh *Sourcehub) AddSource(source any) {
	var ds *deltaSource
	switch v := source.(type) {
	case nmea2000source.Nmea2000Source:
		c := v.Source()
		l := v.Label()
		source := sh.n2kConverter.Convert(l, c)
		ds = &deltaSource{source: source, label: v.Label()}
		sh.inputs = append(sh.inputs, ds)
	default:
		break
	}
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
	}
	return sh.output
}
