package pgn

import (
	"math"

	"github.com/wdantuma/signalk-server-go/converter/base"
)

func NewPgn130306() *PgnBase {
	pgn := NewPgnBase(130306)

	pgn.Fields = append(pgn.Fields,
		base.Field{
			Node:   "environment.wind.speedApparent",
			Source: "windSpeed",
			Filter: func(fields base.InputFields) bool {
				return !fields.Contains("reference") || fields["reference"] == "Apparent"
			},
		},
		base.Field{
			Node:   "environment.wind.speedTrue",
			Source: "windSpeed",
			Filter: func(fields base.InputFields) bool {
				return fields.Contains("reference") && fields["reference"] == "True (boat referenced)"
			},
		},
		base.Field{
			Node:   "environment.wind.speedOverGround",
			Source: "windSpeed",
			Filter: func(fields base.InputFields) bool {
				return fields.Contains("reference") && fields["reference"] == "True (ground referenced to North)"
			},
		},
		base.Field{
			Node: "environment.wind.angleApparent",
			Filter: func(fields base.InputFields) bool {
				return !fields.Contains("reference") || fields["reference"] == "Apparent"
			},
			Value: func(fields base.InputFields) interface{} {
				angle := Float64Value(fields["windAngle"])
				if angle > math.Pi {
					angle = angle - math.Pi*2
				}

				return angle
			},
		},
		base.Field{
			Node:   "environment.wind.angleTrueWater",
			Source: "windAngle",
			Filter: func(fields base.InputFields) bool {
				return fields.Contains("reference") && fields["reference"] == "True (boat referenced)"
			},
		},
		base.Field{
			Node:   "environment.wind.directionTrue",
			Source: "windAngle",
			Filter: func(fields base.InputFields) bool {
				return fields.Contains("reference") && fields["reference"] == "True (ground referenced to North)"
			},
		},
		base.Field{
			Node:   "environment.wind.directionMagnetic",
			Source: "windAngle",
			Filter: func(fields base.InputFields) bool {
				return fields.Contains("reference") && fields["reference"] == "Magnetic (ground referenced to Magnetic North)"
			},
		},
	)

	return pgn
}
