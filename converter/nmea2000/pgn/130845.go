package pgn

import (
	"fmt"

	"github.com/wdantuma/signalk-server-go/converter/base"
	"github.com/wdantuma/signalk-server-go/signalkserver/state"
)

var modeMapping = map[string]string{
	"Day":   "day",
	"Night": "night",
}

var colorMapping = map[string]string{
	"Red":   "red",
	"Green": "green",
	"Blue":  "blue",
	"White": "white",
}

func getDisplayModeObject(state state.ServerState) map[string]interface{} {
	displayMode := make(map[string]interface{})
	value, ok := state.GetStore().Get(fmt.Sprintf("%s/environment.displaymode", state.GetSelf()))
	if ok {
		for k, v := range MapValue(value.Value) {
			displayMode[k] = v
		}
	}
	return displayMode
}

func NewPgn130845() *PgnBase {
	pgn := NewPgnBase(130845)

	pgn.Fields = append(pgn.Fields,
		base.Field{
			Node: "environment.displaymode",
			Filter: func(fields base.InputFields) bool {
				return fields.Contains("key") && fields.Contains("value") && fields["key"] == "Night mode"
			},
			Value: func(fields base.InputFields) interface{} {
				displayMode := getDisplayModeObject(pgn.State)
				displayMode["mode"] = modeMapping[StringValue(fields["value"])]
				return displayMode
			},
		},
		base.Field{
			Node: "environment.displaymode",
			Filter: func(fields base.InputFields) bool {
				return fields.Contains("key") && fields.Contains("value") && fields["key"] == "Backlight level"
			},
			Value: func(fields base.InputFields) interface{} {
				displayMode := getDisplayModeObject(pgn.State)
				displayMode["backlight"] = fields["value"]
				return displayMode
			},
		},
		base.Field{
			Node: "environment.displaymode",
			Filter: func(fields base.InputFields) bool {
				return fields.Contains("key") && fields.Contains("value") && fields["key"] == "Night mode color"
			},
			Value: func(fields base.InputFields) interface{} {
				displayMode := getDisplayModeObject(pgn.State)
				displayMode["color"] = colorMapping[StringValue(fields["value"])]
				return displayMode
			},
		},
		// field{
		// 	node: "environment.displaymode",
		// 	filter: func(fields n2kFields) bool {
		// 		return fields.Contains("key") && fields.Contains("value") && fields["key"] == "Time hour display"
		// 	},
		// 	value: func(fields n2kFields) interface{} {
		// 		displayMode := getDisplayModeObject(pgn.State)
		// 		displayMode["timehour"] = fields["value"]
		// 		return displayMode
		// 	},
		// },
	)

	return pgn
}
