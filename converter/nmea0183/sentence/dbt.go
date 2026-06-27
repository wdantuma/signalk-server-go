package sentence

import "github.com/wdantuma/signalk-server-go/converter/base"

func NewDBT() *SentenceBase {
	pgn := NewSentenceBase("DBT")

	pgn.Fields = append(pgn.Fields,
		base.Field{
			Node:   "environment.depth.belowTransducer",
			Source: "depth",
		},
		// field{
		// 	node:   "environment.depth.surfaceToTransducer",
		// 	source: "offset",
		// },
		// field{
		// 	node: "environment.depth.belowSurface",
		// 	filter: func(fields n2kFields) bool {
		// 		return fields.Contains("depth") && fields.Contains("offset") && Float64Value(fields["offset"]) > 0
		// 	},
		// 	value: func(fields n2kFields) interface{} {
		// 		return Float64Value(fields["depth"]) - Float64Value(fields["offset"])
		// 	},
		// },
		// field{
		// 	node: "environment.depth.belowKeel",
		// 	filter: func(fields n2kFields) bool {
		// 		return fields.Contains("depth") && fields.Contains("offset") && Float64Value(fields["offset"]) < 0
		// 	},
		// 	value: func(fields n2kFields) interface{} {
		// 		return Float64Value(fields["depth"]) + Float64Value(fields["offset"])
		// 	},
		// },
	)

	return pgn
}
