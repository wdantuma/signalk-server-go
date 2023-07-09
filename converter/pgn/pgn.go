package pgn

import (
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
	filter func(n2kFields) bool
	value  func(n2kFields) float64
	node   string
	source string
}

type pgnBase struct {
	Pgn     uint
	PgnInfo *canboat.PGNInfo
	Canboat *canboat.Canboat
	Fields  []field
}

type Pgn interface {
	Convert(can.Frame, *socketcan.CanSource) (signalk.DeltaJson, bool)
}

type PgnBase interface {
	Init(pgn uint, canboat *canboat.Canboat) bool
}

func NewPgnBase() pgnBase {
	return pgnBase{Fields: make([]field, 0)}

}

func (base *pgnBase) GetDelta(frame can.Frame, canSource *socketcan.CanSource) signalk.DeltaJson {
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

func (pgn *pgnBase) Convert(frame can.Frame, canSource *socketcan.CanSource) (signalk.DeltaJson, bool) {
	delta := pgn.GetDelta(frame, canSource)

	fields := make(n2kFields)

	for _, field := range pgn.PgnInfo.Fields.Field {
		value := float64(frame.Data.UnsignedBitsLittleEndian(uint8(field.BitOffset), uint8(field.BitLength))) * float64(field.Resolution)
		fields[field.Id] = value
		if field.FieldType == "LOOKUP" {
			refValue, ok := pgn.GetEnumValueName(value, field.LookupEnumeration)
			if ok {
				fields[field.Id] = refValue
			}
		}
	}

	var include bool = false
	for _, field := range pgn.Fields {

		val := signalk.DeltaJsonUpdatesElemValuesElem{}
		val.Path = field.node
		if field.source != "" {
			value, ok := fields[field.source]
			if !ok {
				log.Println("Source not found")
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

	return delta, include
}

func (base *pgnBase) Init(pgn uint, canboat *canboat.Canboat) bool {
	base.Pgn = pgn
	base.Canboat = canboat
	pgnInfo, ok := canboat.GetPGNInfo(pgn)
	if !ok {
		return false
	}
	base.PgnInfo = pgnInfo
	return true
}

func (base *pgnBase) GetEnumValueName(value float64, name string) (string, bool) {

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
