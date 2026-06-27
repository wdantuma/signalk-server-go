package pgn

import "github.com/wdantuma/signalk-server-go/converter/base"

func NewPgn129039() *PgnBase {
	pgn := NewPgnBase(129039)

	pgn.Fields = append(pgn.Fields,
		base.Field{
			Node:   "navigation.speedOverGround",
			Source: "sog",
		},
		base.Field{
			Node:   "navigation.courseOverGroundTrue",
			Source: "cog",
		},
		base.Field{
			Node: "navigation.position",
			Filter: func(fields base.InputFields) bool {
				return fields.Contains("latitude") && fields.Contains("longitude")
			},
			Value: func(fields base.InputFields) interface{} {
				pos := make(map[string]interface{})
				pos["longitude"] = fields["longitude"]
				pos["latitude"] = fields["latitude"]
				return pos
			},
		},
		base.Field{
			Node:   "navigation.headingTrue",
			Source: "heading",
		},
		base.Field{
			Node: "",
			Filter: func(fields base.InputFields) bool {
				return fields.Contains("userId")
			},
			Value: func(fields base.InputFields) interface{} {
				mmsiReport := make(map[string]interface{})
				mmsiReport["mmsi"] = fields["userId"]
				return mmsiReport
			},
		},
		base.Field{
			Context: GetMmsiContext,
		},
		base.Field{
			Node: "sensors.ais.class",
			Value: func(fields base.InputFields) interface{} {
				return "B"
			},
		},
	)

	return pgn
}
