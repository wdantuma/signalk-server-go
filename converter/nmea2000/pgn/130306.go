package pgn

import (
	"math"
)

func NewPgn130306() *PgnBase {
	pgn := NewPgnBase(130306)

	pgn.Fields = append(pgn.Fields,
		field{
			node:   "environment.wind.speedApparent",
			source: "windSpeed",
			filter: func(fields n2kFields) bool {
				return !fields.Contains("reference") || fields["reference"] == "Apparent"
			},
		},
		field{
			node:   "environment.wind.speedTrue",
			source: "windSpeed",
			filter: func(fields n2kFields) bool {
				return fields.Contains("reference") && fields["reference"] == "True (boat referenced)"
			},
		},
		field{
			node:   "environment.wind.speedOverGround",
			source: "windSpeed",
			filter: func(fields n2kFields) bool {
				return fields.Contains("reference") && fields["reference"] == "True (ground referenced to North)"
			},
		},
		field{
			node: "environment.wind.angleApparent",
			filter: func(fields n2kFields) bool {
				return !fields.Contains("reference") || fields["reference"] == "Apparent"
			},
			value: func(fields n2kFields) interface{} {
				angle := Float64Value(fields["windAngle"])
				if angle > math.Pi {
					angle = angle - math.Pi*2
				}

				return angle
			},
		},
		field{
			node:   "environment.wind.angleTrueWater",
			source: "windAngle",
			filter: func(fields n2kFields) bool {
				return fields.Contains("reference") && fields["reference"] == "True (boat referenced)"
			},
		},
		field{
			node:   "environment.wind.directionTrue",
			source: "windAngle",
			filter: func(fields n2kFields) bool {
				return fields.Contains("reference") && fields["reference"] == "True (ground referenced to North)"
			},
		},
		field{
			node:   "environment.wind.directionMagnetic",
			source: "windAngle",
			filter: func(fields n2kFields) bool {
				return fields.Contains("reference") && fields["reference"] == "Magnetic (ground referenced to Magnetic North)"
			},
		},
	)

	return pgn
}
