package pgn

func NewPgn128267() *PgnBase {
	pgn := NewPgnBase(128267)

	pgn.Fields = append(pgn.Fields,
		field{
			node:   "environment.depth.belowTransducer",
			source: "depth",
		},
		field{
			node:   "environment.depth.surfaceToTransducer",
			source: "offset",
		},
		field{
			node: "environment.depth.belowSurface",
			filter: func(fields n2kFields) bool {
				return fields.Contains("depth") && fields.Contains("offset") && Float64Value(fields["offset"]) > 0
			},
			value: func(fields n2kFields) interface{} {
				return Float64Value(fields["depth"]) - Float64Value(fields["offset"])
			},
		},
		field{
			node: "environment.depth.belowKeel",
			filter: func(fields n2kFields) bool {
				return fields.Contains("depth") && fields.Contains("offset") && Float64Value(fields["offset"]) < 0
			},
			value: func(fields n2kFields) interface{} {
				return Float64Value(fields["depth"]) + Float64Value(fields["offset"])
			},
		},
	)

	return pgn
}
