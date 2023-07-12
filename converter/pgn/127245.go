package pgn

func NewPgn127245() *PgnBase {
	pgn := NewPgnBase(127245)

	pgn.Fields = append(pgn.Fields,
		field{
			node:   "steering.rudderAngle",
			source: "position",
		},
	)

	return pgn
}
