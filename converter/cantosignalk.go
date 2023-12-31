package converter

import (
	"log"

	"github.com/wdantuma/signalk-server-go/canboat"
	"github.com/wdantuma/signalk-server-go/converter/pgn"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/signalkserver/state"
	"github.com/wdantuma/signalk-server-go/source"
)

type canToSignalk struct {
	canboat *canboat.Canboat
	pgn     map[uint]*pgn.PgnBase
	state   state.ServerState
}

type CanToSignalk interface {
	Convert(<-chan source.SourceFrame) <-chan signalk.DeltaJson
}

func NewCanToSignalk(state state.ServerState) (*canToSignalk, error) {
	canboat, err := canboat.NewCanboat()
	if err != nil {
		log.Fatal(err)
	}
	c := canToSignalk{canboat: canboat, state: state, pgn: make(map[uint]*pgn.PgnBase)}
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

func (c *canToSignalk) addPgn(b *pgn.PgnBase) {
	if b.Init(c.canboat, c.state) {
		c.pgn[b.Pgn] = b
	}
}

func (c *canToSignalk) getPgnConverter(frame source.ExtendedFrame) (*pgn.PgnBase, bool) {
	pgn := frame.ID & 0x03FFFF00 >> 8
	pgnConverter, ok := c.pgn[uint(pgn)]
	if ok {
		return pgnConverter, true
	}
	return nil, false
}

func (c *canToSignalk) Convert(canSource <-chan source.SourceFrame) <-chan signalk.DeltaJson {
	output := make(chan signalk.DeltaJson)
	fastframes := make(map[string]*source.ExtendedFrame)
	go func() {
		for {
			sourceFrame, ok := <-canSource
			if ok {
				frame := source.NewExtendedFrame(&sourceFrame.Frame)
				pgnConverter, ok := c.getPgnConverter(frame)
				if ok {
					if pgnConverter.PgnInfo.Type == "Fast" {
						seqNr := frame.UnsignedBitsLittleEndian(0, 4)
						//frameNr := frame.UnsignedBitsLittleEndian(4, 4)
						if seqNr == 0 {
							fastframes[sourceFrame.Label] = frame.First()
							continue
						} else {
							if fastframes[sourceFrame.Label].Next(frame) {
								frame = *fastframes[sourceFrame.Label]
								delete(fastframes, sourceFrame.Label)
							} else {
								continue
							}
						}
					}

					delta, convertOk := pgnConverter.Convert(frame, sourceFrame.Label)
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
