package nmea0183

import (
	"log"

	"github.com/adrianmo/go-nmea"
	"github.com/wdantuma/signalk-server-go/converter/nmea0183/sentence"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/signalkserver/state"
)

type nmea0183ToSignalk struct {
	sentence map[string]*sentence.SentenceBase
	state    state.ServerState
}

type Nmea0183ToSignalk interface {
	Convert(string, <-chan nmea.Sentence) <-chan signalk.DeltaJson
}

func NewNmea0183ToSignalk(state state.ServerState) (*nmea0183ToSignalk, error) {
	c := nmea0183ToSignalk{state: state, sentence: make(map[string]*sentence.SentenceBase)}
	c.addSentence(sentence.NewDBT())
	return &c, nil
}

func (c *nmea0183ToSignalk) addSentence(b *sentence.SentenceBase) {
	c.sentence[b.Sentence] = b
}

func (c *nmea0183ToSignalk) getSentenceConverter(sentence nmea.Sentence) (*sentence.SentenceBase, bool) {
	log.Println(sentence.DataType())
	pgnConverter, ok := c.sentence[sentence.DataType()]
	if ok {
		return pgnConverter, true
	}
	return nil, false
}

func (c *nmea0183ToSignalk) Convert(label string, sentenceSource <-chan nmea.Sentence) <-chan signalk.DeltaJson {
	output := make(chan signalk.DeltaJson)
	go func() {
		for {
			sentence, ok := <-sentenceSource
			if ok {
				sentenceConverter, ok := c.getSentenceConverter(sentence)
				if ok {
					delta, convertOk := sentenceConverter.Convert(sentence, label)
					if convertOk && delta.Context != nil {
						output <- delta
					}
				} else {
					if c.state.GetDebug() {
						log.Printf("Sentence:%s\n", sentence.DataType())
					}
				}
			} else {
				break
			}
		}

		close(output)
	}()

	return output
}
