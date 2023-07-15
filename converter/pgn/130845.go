package pgn

import (
	"fmt"
)

func NewPgn130845() *PgnBase {
	pgn := NewPgnBase(130845)

	pgn.Fields = append(pgn.Fields,
		field{
			node: "environment.displaymode",
			filter: func(fields n2kFields) bool {
				return fields.Contains("key") && fields.Contains("value") && fields["key"] == "Night mode"
			},
			value: func(fields n2kFields) interface{} {
				self := pgn.State.GetSelf()
				displayMode := make(map[string]interface{})
				value, ok := pgn.State.GetStore().Get(fmt.Sprintf("%s/environment.displaymode", self))
				if ok {
					for k, v := range MapValue(value.Value) {
						displayMode[k] = v
					}
				}
				displayMode["mode"] = fields["value"]
				return displayMode
			},
		},
		field{
			node: "environment.displaymode",
			filter: func(fields n2kFields) bool {
				return fields.Contains("key") && fields.Contains("value") && fields["key"] == "Backlight level"
			},
			value: func(fields n2kFields) interface{} {
				self := pgn.State.GetSelf()
				displayMode := make(map[string]interface{})
				value, ok := pgn.State.GetStore().Get(fmt.Sprintf("%s/environment.displaymode", self))
				if ok {
					for k, v := range MapValue(value.Value) {
						displayMode[k] = v
					}
				}
				displayMode["backlight"] = fields["value"]
				return displayMode
			},
		},
		field{
			node: "environment.displaymode",
			filter: func(fields n2kFields) bool {
				return fields.Contains("key") && fields.Contains("value") && fields["key"] == "Night mode color"
			},
			value: func(fields n2kFields) interface{} {
				self := pgn.State.GetSelf()
				displayMode := make(map[string]interface{})
				value, ok := pgn.State.GetStore().Get(fmt.Sprintf("%s/environment.displaymode", self))
				if ok {
					for k, v := range MapValue(value.Value) {
						displayMode[k] = v
					}
				}
				displayMode["color"] = fields["value"]
				return displayMode
			},
		},
		// field{
		// 	node: "environment.displaymode",
		// 	filter: func(fields n2kFields) bool {
		// 		return fields.Contains("key") && fields.Contains("value") && fields["key"] == "Time hour display"
		// 	},
		// 	value: func(fields n2kFields) interface{} {
		// 		self := pgn.State.GetSelf()
		// 		displayMode := make(map[string]interface{})
		// 		value, ok := pgn.State.GetStore().Get(fmt.Sprintf("%s/environment.displaymode", self))
		// 		if ok {
		// 			for k, v := range MapValue(value.Value) {
		// 				displayMode[k] = v
		// 			}
		// 		}
		// 		displayMode["timehour"] = fields["value"]
		// 		return displayMode
		// 	},
		// },
	)

	return pgn
}
