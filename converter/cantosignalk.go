package converter

import (
	"log"
	"strconv"
	"time"

	"github.com/wdantuma/signalk-server-go/canboat"
	"github.com/wdantuma/signalk-server-go/ref"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/signalkserver"
	"github.com/wdantuma/signalk-server-go/socketcan"
)

type canToSignalk struct {
	canboat *canboat.Canboat
}

func NewCanToSignalk() (*canToSignalk, error) {
	canboat, err := canboat.NewCanboat()
	if err != nil {
		log.Fatal(err)
	}
	c := canToSignalk{canboat: canboat}

	return &c, nil
}

func (c *canToSignalk) Convert(canSource *socketcan.CanSource) <-chan signalk.DeltaJson {
	output := make(chan signalk.DeltaJson)
	go func() {
		for frame := range canSource.Source {
			pgn := frame.ID & 0x03FFFF00 >> 8
			pgnInfo, ok := c.canboat.GetPGNInfo(uint(pgn))
			if ok {
				src := frame.ID & 0xFF
				delta := signalk.DeltaJson{}
				delta.Context = ref.String(signalkserver.SELF)
				update := signalk.DeltaJsonUpdatesElem{}
				update.Timestamp = ref.UTCTimeStamp(time.Now()) // TODO get from source
				update.Source = &signalk.Source{Pgn: ref.Float64(float64(pgn)),
					Src:   ref.String(strconv.FormatUint(uint64(src), 10)),
					Type:  "NMEA2000",
					Label: canSource.Label}
				for _, field := range pgnInfo.Fields.Field {
					val := signalk.DeltaJsonUpdatesElemValuesElem{}

					path := ""
					if field.Id == "windSpeed" {
						path = "environment.wind.speedApparent"
					}

					if field.Id == "windAngle" {
						path = "environment.wind.angleApparent"
					}

					val.Path = path
					value := float64(frame.Data.UnsignedBitsLittleEndian(uint8(field.BitOffset), uint8(field.BitLength))) * float64(field.Resolution)
					val.Value = value
					update.Values = append(update.Values, val)
				}
				delta.Updates = append(delta.Updates, update)
				output <- delta
			}
		}
		close(output)
	}()

	return output
}
