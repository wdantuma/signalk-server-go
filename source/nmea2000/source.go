package nmea2000

import "go.einride.tech/can"

type Nmea2000Source interface {
	Source() <-chan can.Frame
	Label() string
}
