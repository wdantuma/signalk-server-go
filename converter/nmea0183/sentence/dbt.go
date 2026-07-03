package sentence

import "github.com/wdantuma/signalk-server-go/converter/base"

func NewDBT() *SentenceBase {
	pgn := NewSentenceBase("DBT")

	pgn.Fields = append(pgn.Fields,
		base.Field{
			Node:   "environment.depth.belowTransducer",
			Source: "DepthMeters",
		},
	)

	return pgn
}
