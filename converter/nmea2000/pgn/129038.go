package pgn

import "github.com/wdantuma/signalk-server-go/converter/base"

var stateMapping = map[string]string{
	"Under way using engine":             "motoring",
	"At anchor":                          "anchored",
	"Not under command":                  "not under command",
	"Restricted manoeuverability":        "restricted manouverability",
	"Constrained by her draught":         "constrained by draft",
	"Moored":                             "moored",
	"Aground":                            "aground",
	"Engaged in Fishing":                 "fishing",
	"Under way sailing":                  "sailing",
	"Hazardous material, High Speed":     "hazardous material high speed",
	"Hazardous material, Wing in Ground": "hazardous material wing in ground",
	"AIS-SART":                           "ais-sart",
}

var specialManeuverMapping = map[string]string{
	"Not available":                   "not available",
	"Not engaged in special maneuver": "not engaged",
	"Engaged in special maneuver":     "engaged",
	"Reserverd":                       "reserved",
}

func NewPgn129038() *PgnBase {
	pgn := NewPgnBase(129038)

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
			Node:   "navigation.rateOfTurn",
			Source: "rateOfTurn",
		},
		base.Field{
			Node:   "navigation.headingTrue",
			Source: "heading",
		},
		base.Field{
			Node: "navigation.state",
			Filter: func(fields base.InputFields) bool {
				return fields.Contains("navStatus")
			},
			Value: func(fields base.InputFields) interface{} {
				return stateMapping[StringValue(fields["navStatus"])]
			},
		},
		base.Field{
			Node: "navigation.specialManeuver",
			Filter: func(fields base.InputFields) bool {
				return fields.Contains("specialManeuverIndicator")
			},
			Value: func(fields base.InputFields) interface{} {
				return specialManeuverMapping[StringValue(fields["specialManeuverIndicator"])]
			},
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
				return "A"
			},
		},
	)

	return pgn
}
