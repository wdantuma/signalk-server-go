package pgn

import "github.com/wdantuma/signalk-server-go/converter/base"

func NewPgn128259() *PgnBase {
	pgn := NewPgnBase(128259)

	pgn.Fields = append(pgn.Fields,
		base.Field{
			Node:   "navigation.speedThroughWater",
			Source: "speedWaterReferenced",
		},
		base.Field{
			Node:   "navigation.speedOverGround",
			Source: "speedGroundReferenced",
		},
		base.Field{
			Node:   "navigation.speedThroughWaterReferenceType",
			Source: "speedWaterReferencedType",
		},
	)

	return pgn
}
