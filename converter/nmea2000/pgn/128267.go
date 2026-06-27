package pgn

import "github.com/wdantuma/signalk-server-go/converter/base"

func NewPgn128267() *PgnBase {
	pgn := NewPgnBase(128267)

	pgn.Fields = append(pgn.Fields,
		base.Field{
			Node:   "environment.depth.belowTransducer",
			Source: "depth",
		},
		base.Field{
			Node:   "environment.depth.surfaceToTransducer",
			Source: "offset",
		},
		base.Field{
			Node: "environment.depth.belowSurface",
			Filter: func(fields base.InputFields) bool {
				return fields.Contains("depth") && fields.Contains("offset") && Float64Value(fields["offset"]) > 0
			},
			Value: func(fields base.InputFields) interface{} {
				return Float64Value(fields["depth"]) - Float64Value(fields["offset"])
			},
		},
		base.Field{
			Node: "environment.depth.belowKeel",
			Filter: func(fields base.InputFields) bool {
				return fields.Contains("depth") && fields.Contains("offset") && Float64Value(fields["offset"]) < 0
			},
			Value: func(fields base.InputFields) interface{} {
				return Float64Value(fields["depth"]) + Float64Value(fields["offset"])
			},
		},
	)

	return pgn
}
