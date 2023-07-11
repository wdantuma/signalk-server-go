package converter

import (
	"log"

	"github.com/wdantuma/signalk-server-go/canboat"
	"github.com/wdantuma/signalk-server-go/converter/pgn"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/signalkserver/state"
	"github.com/wdantuma/signalk-server-go/socketcan"
)

type canToSignalk struct {
	canboat *canboat.Canboat
	pgn     map[uint]*pgn.PgnBase
}

func NewCanToSignalk() (*canToSignalk, error) {
	canboat, err := canboat.NewCanboat()
	if err != nil {
		log.Fatal(err)
	}
	c := canToSignalk{canboat: canboat, pgn: make(map[uint]*pgn.PgnBase)}
	c.addPgn(pgn.NewPgn130306())
	c.addPgn(pgn.NewPgn129038())
	c.addPgn(pgn.NewPgn129039())

	return &c, nil
}

func (c *canToSignalk) addPgn(b *pgn.PgnBase) {
	if b.Init(c.canboat) {
		c.pgn[b.Pgn] = b
	}
}

func (c *canToSignalk) GetPgnConverter(frame socketcan.ExtendedFrame) (*pgn.PgnBase, bool) {
	pgn := frame.ID & 0x03FFFF00 >> 8
	pgnConverter, ok := c.pgn[uint(pgn)]
	if ok {
		return pgnConverter, true
	}
	return nil, false
}

func Reassemble(frame socketcan.ExtendedFrame, length int, input <-chan socketcan.ExtendedFrame) socketcan.ExtendedFrame {
	newBytes := make([]byte, 0)
	newBytes = append(newBytes, frame.Data[2:]...)
	for len(newBytes) < length {
		f := <-input
		newBytes = append(newBytes, f.Data[1:]...)
	}
	frame.Data = newBytes

	return frame
}

func (c *canToSignalk) Convert(state state.ServerState, canSource *socketcan.CanSource) <-chan signalk.DeltaJson {
	output := make(chan signalk.DeltaJson)
	go func() {
		for {
			frame, ok := <-canSource.Source
			if ok {
				pgnConverter, ok := c.GetPgnConverter(frame)
				if ok {
					if pgnConverter.PgnInfo.Type == "Fast" {
						seqNr := frame.UnsignedBitsLittleEndian(0, 4)
						//frameNr := frame.UnsignedBitsLittleEndian(4, 4)
						if seqNr == 0 {
							len := int(frame.UnsignedBitsLittleEndian(8, 8))
							frame = Reassemble(frame, len, canSource.Source)
						}
					}

					delta, convertOk := pgnConverter.Convert(state, frame, canSource)
					if convertOk {
						output <- delta
					}
				} else {
					pgn := frame.ID & 0x03FFFF00 >> 8
					log.Printf("PGN:%d\n", pgn)
				}
			} else {
				break
			}
		}

		close(output)
	}()

	return output
}
