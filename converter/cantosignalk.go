package converter

import (
	"log"

	"github.com/wdantuma/signalk-server-go/canboat"
	"github.com/wdantuma/signalk-server-go/converter/pgn"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/socketcan"
	"go.einride.tech/can"
)

type canToSignalk struct {
	canboat *canboat.Canboat
	pgn     map[uint]interface{}
}

func NewCanToSignalk() (*canToSignalk, error) {
	canboat, err := canboat.NewCanboat()
	if err != nil {
		log.Fatal(err)
	}
	c := canToSignalk{canboat: canboat, pgn: make(map[uint]interface{})}
	c.addPgn(130306, pgn.NewPgn130306())

	return &c, nil
}

func (c *canToSignalk) addPgn(pgn uint, b pgn.PgnBase) {
	if b.Init(pgn, c.canboat) {
		c.pgn[pgn] = b
	}
}

func (c *canToSignalk) GetPgnConverter(frame can.Frame) (interface{}, bool) {
	pgn := frame.ID & 0x03FFFF00 >> 8
	pgnConverter, ok := c.pgn[uint(pgn)]
	if ok {
		return pgnConverter, true
	}
	return nil, false
}

func (c *canToSignalk) Convert(canSource *socketcan.CanSource) <-chan signalk.DeltaJson {
	output := make(chan signalk.DeltaJson)
	go func() {
		for frame := range canSource.Source {
			pgnConverter, ok := c.GetPgnConverter(frame)
			if ok {
				delta, convertOk := pgnConverter.(pgn.Pgn).Convert(frame, canSource)
				if convertOk {
					output <- delta
				}
			}
		}
		close(output)
	}()

	return output
}
