package pgn

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/wdantuma/signalk-server-go/canboat"
	"github.com/wdantuma/signalk-server-go/ref"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/signalkserver"
	"github.com/wdantuma/signalk-server-go/socketcan"
	"go.einride.tech/can"
)

type n2kFields map[string]interface{}

func (field n2kFields) Contains(key string) bool {
	_, ok := field[key]
	return ok
}

type field struct {
	filter  func(n2kFields) bool
	value   func(n2kFields) interface{}
	context func(n2kFields) *string
	node    string
	source  string
}

type PgnBase struct {
	Pgn     uint
	PgnInfo *canboat.PGNInfo
	Canboat *canboat.Canboat
	Fields  []field
}

type Pgn interface {
	Convert(can.Frame, *socketcan.CanSource) (signalk.DeltaJson, bool)
}

func NewPgnBase(pgn uint) *PgnBase {
	return &PgnBase{Pgn: pgn, Fields: make([]field, 0)}

}

func (base *PgnBase) GetDelta(frame socketcan.ExtendedFrame, canSource *socketcan.CanSource) signalk.DeltaJson {
	src := frame.ID & 0xFF
	delta := signalk.DeltaJson{}
	delta.Context = ref.String(signalkserver.SELF)
	update := signalk.DeltaJsonUpdatesElem{}
	update.Timestamp = ref.UTCTimeStamp(time.Now()) // TODO get from source
	update.Source = &signalk.Source{Pgn: ref.Float64(float64(base.Pgn)),
		Src:   ref.String(strconv.FormatUint(uint64(src), 10)),
		Type:  "NMEA2000",
		Label: canSource.Label}

	//update.Values = pgnConverter.Convert(update.Values)
	delta.Updates = append(delta.Updates, update)
	return delta
}

func (pgn *PgnBase) Convert(frame socketcan.ExtendedFrame, canSource *socketcan.CanSource) (signalk.DeltaJson, bool) {
	delta := pgn.GetDelta(frame, canSource)

	fields := make(n2kFields)

	for _, field := range pgn.PgnInfo.Fields.Field {

		switch field.FieldType {
		case "LOOKUP":
			value := float64(frame.UnsignedBitsLittleEndian(int(field.BitOffset), int(field.BitLength))) * float64(field.Resolution)
			if value >= float64(field.RangeMin) && value <= float64(field.RangeMax) {
				refValue, ok := pgn.GetEnumValueName(Float64Value(value), field.LookupEnumeration)
				if ok {
					fields[field.Id] = refValue
				}
			}
		case "NUMBER":
			var value float64
			if field.Signed {
				value = float64(frame.SignedBitsLittleEndian(int(field.BitOffset), int(field.BitLength))) * float64(field.Resolution)
			} else {
				value = float64(frame.UnsignedBitsLittleEndian(int(field.BitOffset), int(field.BitLength))) * float64(field.Resolution)
			}
			if value >= float64(field.RangeMin) && value <= float64(field.RangeMax) {
				fields[field.Id] = value
			}
			break
		case "MMSI":
			var value float64
			value = float64(frame.UnsignedBitsLittleEndian(int(field.BitOffset), int(field.BitLength))) * float64(field.Resolution)
			if value >= float64(field.RangeMin) && value <= float64(field.RangeMax) {
				fields[field.Id] = fmt.Sprintf("%.0f", value)
			}
			break
		}
	}

	var include bool = false
	for _, field := range pgn.Fields {

		val := signalk.DeltaJsonUpdatesElemValuesElem{}
		if field.context != nil {
			delta.Context = field.context(fields)
		} else {
			val.Path = field.node
			if field.source != "" {
				value, ok := fields[field.source]
				if !ok {
					log.Printf("Source  (%s) not found", field.source)
					continue
				}
				val.Value = value
			} else if field.value != nil {
				val.Value = field.value(fields)
			} else {
				log.Println("No value function")
				continue
			}
			if field.filter != nil && field.filter(fields) {
				include = true
				delta.Updates[len(delta.Updates)-1].Values = append(delta.Updates[len(delta.Updates)-1].Values, val)
			} else if field.filter == nil {
				include = true
				delta.Updates[len(delta.Updates)-1].Values = append(delta.Updates[len(delta.Updates)-1].Values, val)
			}
		}
	}

	return delta, include
}

func (base *PgnBase) Init(canboat *canboat.Canboat) bool {
	base.Canboat = canboat
	pgnInfo, ok := canboat.GetPGNInfo(base.Pgn)
	if !ok {
		return false
	}
	base.PgnInfo = pgnInfo
	return true
}

func (base *PgnBase) GetEnumValueName(value float64, name string) (string, bool) {

	lookupEnumeration, ok := base.Canboat.GetLookupEnumeration(name)
	if ok {
		for _, v := range lookupEnumeration.EnumPair {
			if v.ValueAttr == uint(value) {
				return v.Name, true
			}
		}
	}

	return "", false
}

func Float64Value(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	default:
		return 0
	}
}

func StringValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	default:
		return ""
	}
}

func GetMmsiContext(fields n2kFields) *string {
	if fields.Contains("userId") {
		mmsi := fmt.Sprintf("vessels.urn:mrn:imo:mmsi:%s", StringValue(fields["userId"]))
		return &mmsi
	} else {
		return nil
	}
}
