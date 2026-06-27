package pgn

import "github.com/wdantuma/signalk-server-go/converter/base"

func NewPgn65359() *PgnBase {
	pgn := NewPgnBase(65359)

	pgn.Fields = append(pgn.Fields,
		base.Field{
			Node:   "navigation.headingTrue",
			Source: "headingTrue",
		},
		base.Field{
			Node:   "navigation.headingMagnetic",
			Source: "headingMagnetic",
		},
	)

	return pgn
}
