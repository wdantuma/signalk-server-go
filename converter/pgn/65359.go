package pgn

func NewPgn65359() *PgnBase {
	pgn := NewPgnBase(65359)

	pgn.Fields = append(pgn.Fields,
		field{
			node:   "navigation.headingTrue",
			source: "headingTrue",
		},
		field{
			node:   "navigation.headingMagnetic",
			source: "headingMagnetic",
		},
	)

	return pgn
}
