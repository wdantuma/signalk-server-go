package pgn

import (
	"math"
)

type Pgn130306 struct {
	pgnBase
}

func NewPgn130306() *Pgn130306 {
	pgn := &Pgn130306{pgnBase: NewPgnBase()}

	pgn.Fields = append(pgn.Fields,
		field{
			node:   "environment.wind.speedApparent",
			source: "windSpeed",
			filter: func(fields n2kFields) bool {
				return !fields.Contains("reference") || fields["reference"] == "Apparent"
			},
		},
		field{
			node: "environment.wind.angleApparent",
			filter: func(fields n2kFields) bool {
				return !fields.Contains("reference") || fields["reference"] == "Apparent"
			},
			value: func(fields n2kFields) float64 {
				angle := Float64Value(fields["windAngle"])
				if angle > math.Pi {
					angle = angle - math.Pi*2
				}

				return angle
			},
		},
	)

	return pgn
}
