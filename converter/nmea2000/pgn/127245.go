package pgn

import "github.com/wdantuma/signalk-server-go/converter/base"

func NewPgn127245() *PgnBase {
	pgn := NewPgnBase(127245)

	pgn.Fields = append(pgn.Fields,
		base.Field{
			Node:   "steering.rudderAngle",
			Source: "position",
		},
	)

	return pgn
}
