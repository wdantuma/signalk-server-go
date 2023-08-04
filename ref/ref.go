package ref

import (
	"time"

	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/signalkserver/state"
)

func String(v interface{}) *string {
	vv := v.(string)
	return &vv
}

func UTCTimeStamp(i time.Time) *signalk.Timestamp {
	return (*signalk.Timestamp)(String(i.UTC().Format(state.TIME_FORMAT)))
}

func Float64(v interface{}) *float64 {
	vv := v.(float64)
	return &vv
}

func Float32(v interface{}) *float32 {
	vv := v.(float32)
	return &vv
}

func Int64(v interface{}) *int64 {
	vv := v.(int64)
	return &vv
}

func Int32(v interface{}) *int32 {
	vv := v.(int32)
	return &vv
}
