package nmea2000

import (
	"go.einride.tech/can"
)

type SourceFrame struct {
	Frame can.Frame
	Label string
}
