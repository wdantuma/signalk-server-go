package pgn

func NewPgn129039() *PgnBase {
	pgn := NewPgnBase(129039)

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
			node:   "navigation.headingTrue",
			source: "heading",
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
				return "B"
			},
		},
	)

	return pgn
}
