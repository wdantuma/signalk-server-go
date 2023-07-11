package pgn

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
		field{
			node:   "navigation.speedOverGround",
			source: "sog",
		},
		field{
			node:   "navigation.courseOverGroundTrue",
			source: "cog",
		},
		field{
			node: "navigation.position",
			filter: func(fields n2kFields) bool {
				return fields.Contains("latitude") && fields.Contains("longitude")
			},
			value: func(fields n2kFields) interface{} {
				pos := make(map[string]interface{})
				pos["longitude"] = fields["longitude"]
				pos["latitude"] = fields["latitude"]
				return pos
			},
		},
		field{
			node:   "navigation.rateOfTurn",
			source: "rateOfTurn",
		},
		field{
			node:   "navigation.headingTrue",
			source: "heading",
		},
		field{
			node: "navigation.state",
			filter: func(fields n2kFields) bool {
				return fields.Contains("navStatus")
			},
			value: func(fields n2kFields) interface{} {
				return stateMapping[StringValue(fields["navStatus"])]
			},
		},
		field{
			node: "navigation.specialManeuver",
			filter: func(fields n2kFields) bool {
				return fields.Contains("specialManeuverIndicator")
			},
			value: func(fields n2kFields) interface{} {
				return specialManeuverMapping[StringValue(fields["specialManeuverIndicator"])]
			},
		},
		field{
			node: "",
			filter: func(fields n2kFields) bool {
				return fields.Contains("userId")
			},
			value: func(fields n2kFields) interface{} {
				mmsiReport := make(map[string]interface{})
				mmsiReport["mmsi"] = fields["userId"]
				return mmsiReport
			},
		},
		field{
			context: GetMmsiContext,
		},
		field{
			node: "sensors.ais.class",
			value: func(fields n2kFields) interface{} {
				return "A"
			},
		},
	)

	return pgn
}
