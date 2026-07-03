package sentence

import (
	"log"
	"reflect"
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/wdantuma/signalk-server-go/converter/base"
	"github.com/wdantuma/signalk-server-go/ref"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/signalkserver/state"
)

type SentenceBase struct {
	Sentence string
	Fields   []base.Field
	State    state.ServerState
}

func NewSentenceBase(sentence string) *SentenceBase {
	return &SentenceBase{Sentence: sentence, Fields: make([]base.Field, 0)}

}

func (baseSentence *SentenceBase) getDelta(sentence nmea.Sentence, source string) signalk.DeltaJson {
	delta := signalk.DeltaJson{}
	delta.Context = ref.String(baseSentence.State.GetSelf())
	update := signalk.DeltaJsonUpdatesElem{}
	update.Timestamp = ref.UTCTimeStamp(time.Now()) // TODO get from source
	update.Source = &signalk.Source{
		Sentence: &baseSentence.Sentence,
		Type:     "NMEA0183",
		Label:    source,
	}

	//update.Values = pgnConverter.Convert(update.Values)
	delta.Updates = append(delta.Updates, update)
	return delta
}

func (baseSentence *SentenceBase) Convert(sentence nmea.Sentence, source string) (signalk.DeltaJson, bool) {
	delta := baseSentence.getDelta(sentence, source)
	fields := make(base.InputFields)
	metadata := make(map[string]signalk.Meta)

	r := reflect.ValueOf(sentence)
	for i := 0; i < r.NumField(); i++ {
		fields[r.Type().Field(i).Name] = r.Field(i).Interface()
	}
	var include bool = false
	for _, field := range baseSentence.Fields {

		val := signalk.DeltaJsonUpdatesElemValuesElem{}
		meta := signalk.DeltaJsonUpdatesElemMetaElem{}
		if field.Context != nil {
			delta.Context = field.Context(fields)
		} else {
			val.Path = field.Node
			meta.Path = field.Node
			if field.Source != "" {
				value, ok := fields[field.Source]
				if !ok {
					continue
				}
				m, ok := metadata[field.Source]
				if ok {
					meta.Value = m
				}
				val.Value = value
			} else if field.Value != nil {
				val.Value = field.Value(fields)
			} else {
				log.Println("No value function")
				continue
			}
			if (field.Filter != nil && field.Filter(fields)) || field.Filter == nil {
				include = true
				delta.Updates[len(delta.Updates)-1].Meta = append(delta.Updates[len(delta.Updates)-1].Meta, meta)
				delta.Updates[len(delta.Updates)-1].Values = append(delta.Updates[len(delta.Updates)-1].Values, val)
			}
		}
	}

	return delta, include
}

func (baseSentence *SentenceBase) Init(state state.ServerState) bool {
	baseSentence.State = state
	return true
}

func Float64Value(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	default:
		return 0
	}
}
