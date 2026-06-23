package nmea0183

type Nmea0183Source interface {
	Source() chan Sentence
	Label() string
}
