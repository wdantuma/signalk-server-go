package pgn

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/wdantuma/signalk-server-go/canboat"
	"github.com/wdantuma/signalk-server-go/converter/base"
	"github.com/wdantuma/signalk-server-go/ref"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/signalkserver/state"
	"github.com/wdantuma/signalk-server-go/source/nmea2000"
)

type PgnBase struct {
	Pgn     uint
	PgnInfo *canboat.PGNInfo
	Canboat *canboat.Canboat
	Fields  []base.Field
	State   state.ServerState
}

func NewPgnBase(pgn uint) *PgnBase {
	return &PgnBase{Pgn: pgn, Fields: make([]base.Field, 0)}

}

func (base *PgnBase) getDelta(frame nmea2000.ExtendedFrame, source string) signalk.DeltaJson {
	src := frame.ID & 0xFF
	delta := signalk.DeltaJson{}
	delta.Context = ref.String(base.State.GetSelf())
	update := signalk.DeltaJsonUpdatesElem{}
	update.Timestamp = ref.UTCTimeStamp(time.Now()) // TODO get from source
	update.Source = &signalk.Source{
		Pgn:   ref.Float64(float64(base.Pgn)),
		Src:   ref.String(strconv.FormatUint(uint64(src), 10)),
		Type:  "NMEA2000",
		Label: source,
	}

	//update.Values = pgnConverter.Convert(update.Values)
	delta.Updates = append(delta.Updates, update)
	return delta
}

func (pgn *PgnBase) Convert(frame nmea2000.ExtendedFrame, source string) (signalk.DeltaJson, bool) {
	delta := pgn.getDelta(frame, source)

	lookupFieldTypeField := canboat.Field{}

	fields := make(base.InputFields)
	metadata := make(map[string]signalk.Meta)

	for _, f := range pgn.PgnInfo.Fields.Field {

		field := f // copy

		meta := signalk.Meta{}
		unit := f.Unit
		meta.Units = &unit
		meta.Description = f.Description
		metadata[f.Id] = meta

		if field.BitOffset == 0 && field.BitLength == 0 {
			field.BitOffset = lookupFieldTypeField.BitOffset
			field.BitLength = lookupFieldTypeField.BitLength
			field.FieldType = lookupFieldTypeField.FieldType
			field.Unit = lookupFieldTypeField.Unit
			field.Signed = lookupFieldTypeField.Signed
			field.Resolution = lookupFieldTypeField.Resolution
			field.RangeMax = lookupFieldTypeField.RangeMax
			field.RangeMin = lookupFieldTypeField.RangeMin
			field.LookupEnumeration = lookupFieldTypeField.LookupEnumeration
			field.LookupBitEnumeration = lookupFieldTypeField.LookupBitEnumeration
		} else {
			lookupFieldTypeField.BitOffset = field.BitOffset + field.BitLength
		}

		switch field.FieldType {
		case "LOOKUP":
			value := float64(frame.UnsignedBitsLittleEndian(int(field.BitOffset), int(field.BitLength))) * float64(field.Resolution)
			if value >= float64(field.RangeMin) && value <= float64(field.RangeMax) {
				refValue, ok := pgn.Canboat.GetLookupEnumeration(field.LookupEnumeration, Float64Value(value))
				if ok {
					fields[field.Id] = refValue
				}

				fieldType, ok := pgn.Canboat.GetLookupFieldTypeEnumeration(field.LookupFieldTypeEnumeration, Float64Value(value))
				if ok {
					fields[field.Id] = fieldType.Name
					lookupFieldTypeField.FieldType = fieldType.FieldType
					lookupFieldTypeField.Signed = fieldType.Signed
					lookupFieldTypeField.Unit = fieldType.Unit
					lookupFieldTypeField.Resolution = field.Resolution
					lookupFieldTypeField.BitLength = fieldType.Bits
					lookupFieldTypeField.RangeMax = 255 // TODO Fix this
					if fieldType.FieldType == "LOOKUP" {
						lookupFieldTypeField.LookupEnumeration = fieldType.LookupEnumeration
					}
				}

			}
		case "NUMBER":
			var value float64
			if field.Signed {
				value = float64(frame.SignedBitsLittleEndian(int(field.BitOffset), int(field.BitLength))) * float64(field.Resolution)
			} else {
				value = float64(frame.UnsignedBitsLittleEndian(int(field.BitOffset), int(field.BitLength))) * float64(field.Resolution)
			}
			if pgn.State.GetDebug() {
				// do not filter out of limit values
				fields[field.Id] = value
			} else {
				if value >= float64(field.RangeMin) && value <= float64(field.RangeMax) {
					fields[field.Id] = value
				}
			}
		case "MMSI":
			var value float64
			value = float64(frame.UnsignedBitsLittleEndian(int(field.BitOffset), int(field.BitLength))) * float64(field.Resolution)
			if value >= float64(field.RangeMin) && value <= float64(field.RangeMax) {
				fields[field.Id] = fmt.Sprintf("%.0f", value)
			}
		}
	}

	var include bool = false
	for _, field := range pgn.Fields {

		val := signalk.DeltaJsonUpdatesElemValuesElem{}
		meta := signalk.DeltaJsonUpdatesElemMetaElem{}
		if field.Context != nil {
			delta.Context = field.Context(fields)
		} else {
			val.Path = field.Node
			meta.Path = field.Node
			if field.Source != "" {
				value, ok := fields[field.Source]
				if !ok {
					continue
				}
				m, ok := metadata[field.Source]
				if ok {
					meta.Value = m
				}
				val.Value = value
			} else if field.Value != nil {
				val.Value = field.Value(fields)
			} else {
				log.Println("No value function")
				continue
			}
			if (field.Filter != nil && field.Filter(fields)) || field.Filter == nil {
				include = true
				delta.Updates[len(delta.Updates)-1].Meta = append(delta.Updates[len(delta.Updates)-1].Meta, meta)
				delta.Updates[len(delta.Updates)-1].Values = append(delta.Updates[len(delta.Updates)-1].Values, val)
			}
		}
	}

	return delta, include
}

func (base *PgnBase) Init(canboat *canboat.Canboat, state state.ServerState) bool {
	base.Canboat = canboat
	base.State = state
	pgnInfo, ok := canboat.GetPGNInfo(base.Pgn)
	if !ok {
		return false
	}
	base.PgnInfo = pgnInfo
	return true
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

func MapValue(value interface{}) map[string]interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		return v
	default:
		return nil
	}
}

func GetMmsiContext(fields base.InputFields) *string {
	if fields.Contains("userId") {
		mmsi := fmt.Sprintf("vessels.urn:mrn:imo:mmsi:%s", StringValue(fields["userId"]))
		return &mmsi
	} else {
		return nil
	}
}
