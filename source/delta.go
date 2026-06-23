package source

import "github.com/wdantuma/signalk-server-go/signalk"

type Delta struct {
	Delta signalk.DeltaJson
	Label string
}
