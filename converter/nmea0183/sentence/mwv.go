package sentence

import (
	"math"

	"github.com/wdantuma/signalk-server-go/converter/base"
)

func NewMWV() *SentenceBase {
	pgn := NewSentenceBase("MWV")

	pgn.Fields = append(pgn.Fields,
		base.Field{
			Node: "environment.wind.speedApparent",
			Filter: func(fields base.InputFields) bool {
				return fields["Reference"] == "R"
			},
			Value: func(fields base.InputFields) any {
				switch fields["WindSpeedUnit"] {
				case "N": // knots
					return Float64Value(fields["WindSpeed"]) * 0.51444
				case "M": // Meter / second
					return Float64Value(fields["WindSpeed"])
				default: // KM / hour
					return Float64Value(fields["WindSpeed"]) / 3.6
				}
			},
		},
		base.Field{
			Node: "environment.wind.speedTrue",
			Filter: func(fields base.InputFields) bool {
				return fields["Reference"] == "T"
			},
			Value: func(fields base.InputFields) any {
				switch fields["WindSpeedUnit"] {
				case "N": // knots
					return Float64Value(fields["WindSpeed"]) * 0.51444
				case "M": // Meter / second
					return Float64Value(fields["WindSpeed"])
				default: // KM / hour
					return Float64Value(fields["WindSpeed"]) / 3.6
				}
			},
		},
		base.Field{
			Node: "environment.wind.angleApparent",
			Filter: func(fields base.InputFields) bool {
				return fields["Reference"] == "R"
			},
			Value: func(fields base.InputFields) interface{} {
				angle := Float64Value(fields["WindAngle"]) * math.Pi / 180
				if angle > math.Pi {
					angle = angle - math.Pi*2
				}

				return angle
			},
		},
		base.Field{
			Node: "environment.wind.angleTrue",
			Filter: func(fields base.InputFields) bool {
				return fields["Reference"] == "T"
			},
			Value: func(fields base.InputFields) interface{} {
				angle := Float64Value(fields["WindAngle"]) * math.Pi / 180
				if angle > math.Pi {
					angle = angle - math.Pi*2
				}

				return angle
			},
		},
	)

	return pgn
}
