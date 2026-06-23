package nmea2000

import (
	"log"

	"github.com/wdantuma/signalk-server-go/canboat"
	"github.com/wdantuma/signalk-server-go/converter/nmea2000/pgn"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/signalkserver/state"
	"github.com/wdantuma/signalk-server-go/source/nmea2000"
	"go.einride.tech/can"
)

type nme2000ToSignalk struct {
	canboat *canboat.Canboat
	pgn     map[uint]*pgn.PgnBase
	state   state.ServerState
}

type Nmea2000ToSignalk interface {
	Convert(string, <-chan can.Frame) <-chan signalk.DeltaJson
}

func NewNmea2000ToSignalk(state state.ServerState) (*nme2000ToSignalk, error) {
	canboat, err := canboat.NewCanboat()
	if err != nil {
		log.Fatal(err)
	}
	c := nme2000ToSignalk{canboat: canboat, state: state, pgn: make(map[uint]*pgn.PgnBase)}
	c.addPgn(pgn.NewPgn130306())
	c.addPgn(pgn.NewPgn129038())
	c.addPgn(pgn.NewPgn129039())
	c.addPgn(pgn.NewPgn127245())
	c.addPgn(pgn.NewPgn128267())
	c.addPgn(pgn.NewPgn128259())
	c.addPgn(pgn.NewPgn130845())
	c.addPgn(pgn.NewPgn65359())

	return &c, nil
}

func (c *nme2000ToSignalk) addPgn(b *pgn.PgnBase) {
	if b.Init(c.canboat, c.state) {
		c.pgn[b.Pgn] = b
	}
}

func (c *nme2000ToSignalk) getPgnConverter(frame nmea2000.ExtendedFrame) (*pgn.PgnBase, bool) {
	pgn := frame.ID & 0x03FFFF00 >> 8
	pgnConverter, ok := c.pgn[uint(pgn)]
	if ok {
		return pgnConverter, true
	}
	return nil, false
}

func (c *nme2000ToSignalk) Convert(label string, canSource <-chan can.Frame) <-chan signalk.DeltaJson {
	output := make(chan signalk.DeltaJson)
	fastframes := make(map[string]*nmea2000.ExtendedFrame)
	go func() {
		for {
			canFrame, ok := <-canSource
			if ok {
				frame := nmea2000.NewExtendedFrame(&canFrame)
				pgnConverter, ok := c.getPgnConverter(frame)
				if ok {
					if pgnConverter.PgnInfo.Type == "Fast" {
						seqNr := frame.UnsignedBitsLittleEndian(0, 4)
						//frameNr := frame.UnsignedBitsLittleEndian(4, 4)
						if seqNr == 0 {
							fastframes[label] = frame.First()
							continue
						} else {
							if fastframes[label].Next(frame) {
								frame = *fastframes[label]
								delete(fastframes, label)
							} else {
								continue
							}
						}
					}

					delta, convertOk := pgnConverter.Convert(frame, label)
					if convertOk && delta.Context != nil {
						output <- delta
					}
				} else {
					pgn := frame.ID & 0x03FFFF00 >> 8
					if c.state.GetDebug() {
						log.Printf("PGN:%d\n", pgn)
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
