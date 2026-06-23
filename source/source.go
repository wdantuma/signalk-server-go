package source

import "github.com/wdantuma/signalk-server-go/signalk"

type DeltaSource interface {
	Source() <-chan signalk.DeltaJson
	Label() string
}

type deltaSource struct {
	source <-chan signalk.DeltaJson
	label  string
}

func (ds *deltaSource) Source() <-chan signalk.DeltaJson {
	return ds.source
}

func (ds *deltaSource) Label() string {
	return ds.label
}
