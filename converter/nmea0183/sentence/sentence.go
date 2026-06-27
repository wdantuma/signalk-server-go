package sentence

import (
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/wdantuma/signalk-server-go/ref"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/signalkserver/state"
)

type field struct {
	node   string
	source string
}

type SentenceBase struct {
	Sentence string
	Fields   []field
	State    state.ServerState
}

func NewSentenceBase(sentence string) *SentenceBase {
	return &SentenceBase{Sentence: sentence, Fields: make([]field, 0)}

}

func (base *SentenceBase) getDelta(sentence nmea.Sentence, source string) signalk.DeltaJson {
	delta := signalk.DeltaJson{}
	delta.Context = ref.String(base.State.GetSelf())
	update := signalk.DeltaJsonUpdatesElem{}
	update.Timestamp = ref.UTCTimeStamp(time.Now()) // TODO get from source
	update.Source = &signalk.Source{
		Sentence: &base.Sentence,
		Type:     "NMEA0183",
		Label:    source,
	}

	//update.Values = pgnConverter.Convert(update.Values)
	delta.Updates = append(delta.Updates, update)
	return delta
}

func (base *SentenceBase) Convert(sentence nmea.Sentence, source string) (signalk.DeltaJson, bool) {
	return base.getDelta(sentence, source), true
}
