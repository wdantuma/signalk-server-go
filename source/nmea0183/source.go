package nmea0183

import "github.com/adrianmo/go-nmea"

type Nmea0183Source interface {
	Source() <-chan nmea.Sentence
	Label() string
}
