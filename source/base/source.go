package base

import "github.com/wdantuma/signalk-server-go/signalk"

type DeltaSource interface {
	Source() <-chan signalk.DeltaJson
	Label() string
	Start()
	Stop()
}
