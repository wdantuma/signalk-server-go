package ref

import (
	"time"

	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/signalkserver/state"
)

func String(i string) *string {
	return &i
}

func UTCTimeStamp(i time.Time) *signalk.Timestamp {
	return (*signalk.Timestamp)(String(i.UTC().Format(state.TIME_FORMAT)))
}

func Float64(f float64) *float64 {
	return &f
}
