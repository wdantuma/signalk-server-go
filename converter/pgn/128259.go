package pgn

func NewPgn128259() *PgnBase {
	pgn := NewPgnBase(128259)

	pgn.Fields = append(pgn.Fields,
		field{
			node:   "navigation.speedThroughWater",
			source: "speedWaterReferenced",
		},
		field{
			node:   "navigation.speedOverGround",
			source: "speedGroundReferenced",
		},
		field{
			node:   "navigation.speedThroughWaterReferenceType",
			source: "speedWaterReferencedType",
		},
	)

	return pgn
}
