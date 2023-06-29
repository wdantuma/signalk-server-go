// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package signalk

import "encoding/json"
import "fmt"
import "reflect"

// Maritime Mobile Service Identity (MMSI) for aircraft. Has to be 9 digits. See
// http://en.wikipedia.org/wiki/Maritime_Mobile_Service_Identity for information.
type AircraftMmsi string

type AlarmMethodEnum string

const AlarmMethodEnumSound AlarmMethodEnum = "sound"
const AlarmMethodEnumVisual AlarmMethodEnum = "visual"

type AlarmState string

const AlarmStateAlarm AlarmState = "alarm"
const AlarmStateAlert AlarmState = "alert"
const AlarmStateEmergency AlarmState = "emergency"
const AlarmStateNominal AlarmState = "nominal"
const AlarmStateNormal AlarmState = "normal"
const AlarmStateWarn AlarmState = "warn"

// Maritime Mobile Service Identity (MMSI) for . Has to be 9 digits. See
// http://en.wikipedia.org/wiki/Maritime_Mobile_Service_Identity for information.
type AtonMmsi string

// Filesystem specific data, e.g. security, possibly more later.
type Attr struct {
	// The group owning this resource.
	Group string `json:"_group,omitempty" yaml:"_group,omitempty" mapstructure:"_group,omitempty"`

	// Unix style permissions, often written in `owner:group:other` form, `-rw-r--r--`
	Mode int `json:"_mode,omitempty" yaml:"_mode,omitempty" mapstructure:"_mode,omitempty"`

	// The owner of this resource.
	Owner string `json:"_owner,omitempty" yaml:"_owner,omitempty" mapstructure:"_owner,omitempty"`
}

type CommonValueFields struct {
	// Source corresponds to the JSON schema field "$source".
	Source SourceRef `json:"$source" yaml:"$source" mapstructure:"$source"`

	// Attr corresponds to the JSON schema field "_attr".
	Attr *Attr `json:"_attr,omitempty" yaml:"_attr,omitempty" mapstructure:"_attr,omitempty"`

	// Meta corresponds to the JSON schema field "meta".
	Meta *Meta `json:"meta,omitempty" yaml:"meta,omitempty" mapstructure:"meta,omitempty"`

	// Pgn corresponds to the JSON schema field "pgn".
	Pgn *float64 `json:"pgn,omitempty" yaml:"pgn,omitempty" mapstructure:"pgn,omitempty"`

	// Sentence corresponds to the JSON schema field "sentence".
	Sentence *string `json:"sentence,omitempty" yaml:"sentence,omitempty" mapstructure:"sentence,omitempty"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp Timestamp `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`
}

// Data should be of type number.
type DatetimeValue map[string]interface{}

// Reusable definitions of core object types
type DefinitionsJson map[string]interface{}

// Schema for defining updates and subscriptions to parts of a SignalK data model,
// for example for communicating updates of data
type DeltaJson struct {
	// The context path of the updates, eg. the top level path plus object identifier.
	Context *string `json:"context,omitempty" yaml:"context,omitempty" mapstructure:"context,omitempty"`

	// The updates
	Updates []DeltaJsonUpdatesElem `json:"updates" yaml:"updates" mapstructure:"updates"`
}

type DeltaJsonUpdatesElem struct {
	// Source corresponds to the JSON schema field "$source".
	SourceRef *SourceRef `json:"$source,omitempty" yaml:"$source,omitempty" mapstructure:"$source,omitempty"`

	// Meta corresponds to the JSON schema field "meta".
	Meta []DeltaJsonUpdatesElemMetaElem `json:"meta,omitempty" yaml:"meta,omitempty" mapstructure:"meta,omitempty"`

	// Source_2 corresponds to the JSON schema field "source".
	Source *Source `json:"source,omitempty" yaml:"source,omitempty" mapstructure:"source,omitempty"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp *Timestamp `json:"timestamp,omitempty" yaml:"timestamp,omitempty" mapstructure:"timestamp,omitempty"`

	// Values corresponds to the JSON schema field "values".
	Values []DeltaJsonUpdatesElemValuesElem `json:"values,omitempty" yaml:"values,omitempty" mapstructure:"values,omitempty"`
}

type DeltaJsonUpdatesElemMetaElem struct {
	// The local path to the data value
	Path string `json:"path" yaml:"path" mapstructure:"path"`

	// Value corresponds to the JSON schema field "value".
	Value Meta `json:"value" yaml:"value" mapstructure:"value"`
}

type DeltaJsonUpdatesElemValuesElem struct {
	// The local path to the data value
	Path string `json:"path" yaml:"path" mapstructure:"path"`

	// Value corresponds to the JSON schema field "value".
	Value interface{} `json:"value" yaml:"value" mapstructure:"value"`
}

// A geohash (see http://geohash.org)
type Geohash string

// Schema for defining the hello message passed from the server to a client
// following succesful websocket connection
type HelloJson struct {
	// The name of the Signal K server software
	Name *string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`

	// Playback rate for history playback connections: 1 is real time, 2 is two times
	// and 0.5 half the real time rate
	PlaybackRate *float64 `json:"playbackRate,omitempty" yaml:"playbackRate,omitempty" mapstructure:"playbackRate,omitempty"`

	// The designated roles of the server
	Roles []HelloJsonRolesElem `json:"roles" yaml:"roles" mapstructure:"roles"`

	// This holds the context (prefix + UUID, MMSI or URL in dot notation) of the
	// server's self object.
	Self *string `json:"self,omitempty" yaml:"self,omitempty" mapstructure:"self,omitempty"`

	// Starttime for history playback connections
	StartTime *Timestamp `json:"startTime,omitempty" yaml:"startTime,omitempty" mapstructure:"startTime,omitempty"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp *Timestamp `json:"timestamp,omitempty" yaml:"timestamp,omitempty" mapstructure:"timestamp,omitempty"`

	// Version of the schema and APIs that this data is using in canonical format i.e.
	// 1.5.0.
	Version Version `json:"version" yaml:"version" mapstructure:"version"`
}

type HelloJsonRolesElem string

const HelloJsonRolesElemAux HelloJsonRolesElem = "aux"
const HelloJsonRolesElemMain HelloJsonRolesElem = "main"
const HelloJsonRolesElemMaster HelloJsonRolesElem = "master"
const HelloJsonRolesElemSlave HelloJsonRolesElem = "slave"

// Provides meta data to enable alarm and display configuration.
type Meta struct {
	// The method to use to raise the alarm. An alarm requires immediate attention, eg
	// no oil pressure
	AlarmMethod []MetaAlarmMethodElem `json:"alarmMethod,omitempty" yaml:"alarmMethod,omitempty" mapstructure:"alarmMethod,omitempty"`

	// The method to use to raise the alert. An alert is an event that should be known
	AlertMethod []MetaAlertMethodElem `json:"alertMethod,omitempty" yaml:"alertMethod,omitempty" mapstructure:"alertMethod,omitempty"`

	// Description of the SK path.
	Description string `json:"description" yaml:"description" mapstructure:"description"`

	// A display name for this value. This is shown on the gauge and should not
	// include units.
	DisplayName *string `json:"displayName,omitempty" yaml:"displayName,omitempty" mapstructure:"displayName,omitempty"`

	// Gives details of the display scale against which the value should be displayed
	DisplayScale *MetaDisplayScale `json:"displayScale,omitempty" yaml:"displayScale,omitempty" mapstructure:"displayScale,omitempty"`

	// The method to use to raise an emergency. An emergency is an immediate danger to
	// life or vessel
	EmergencyMethod []MetaEmergencyMethodElem `json:"emergencyMethod,omitempty" yaml:"emergencyMethod,omitempty" mapstructure:"emergencyMethod,omitempty"`

	// List of permissible values
	Enum []interface{} `json:"enum,omitempty" yaml:"enum,omitempty" mapstructure:"enum,omitempty"`

	// gaugeType is deprecated. The type of gauge necessary to display this value.
	GaugeType *string `json:"gaugeType,omitempty" yaml:"gaugeType,omitempty" mapstructure:"gaugeType,omitempty"`

	// A long name for this value.
	LongName *string `json:"longName,omitempty" yaml:"longName,omitempty" mapstructure:"longName,omitempty"`

	// Properties corresponds to the JSON schema field "properties".
	Properties MetaProperties `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// A short name for this value.
	ShortName *string `json:"shortName,omitempty" yaml:"shortName,omitempty" mapstructure:"shortName,omitempty"`

	// The timeout in (fractional) seconds after which this data is invalid.
	Timeout *float64 `json:"timeout,omitempty" yaml:"timeout,omitempty" mapstructure:"timeout,omitempty"`

	// The (derived) SI unit of this value.
	Units *string `json:"units,omitempty" yaml:"units,omitempty" mapstructure:"units,omitempty"`

	// The method to use to raise the warning. A warning is an unexpected event that
	// may require attention
	WarnMethod []MetaWarnMethodElem `json:"warnMethod,omitempty" yaml:"warnMethod,omitempty" mapstructure:"warnMethod,omitempty"`

	// The zones defining the range of values for this signalk value.
	Zones []MetaZonesElem `json:"zones,omitempty" yaml:"zones,omitempty" mapstructure:"zones,omitempty"`
}

type MetaAlarmMethodElem interface{}

type MetaAlertMethodElem interface{}

// Gives details of the display scale against which the value should be displayed
type MetaDisplayScale struct {
	// The suggested lower limit for the pointer (or equivalent) on the display
	Lower *float64 `json:"lower,omitempty" yaml:"lower,omitempty" mapstructure:"lower,omitempty"`

	// The power to use when the displayScale/type is set to 'power'. Can be any
	// numeric value except zero.
	Power *float64 `json:"power,omitempty" yaml:"power,omitempty" mapstructure:"power,omitempty"`

	// The suggested type of scale to use
	Type *MetaDisplayScaleType `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`

	// The suggested upper limit for the pointer (or equivalent) on the display
	Upper *float64 `json:"upper,omitempty" yaml:"upper,omitempty" mapstructure:"upper,omitempty"`
}

type MetaDisplayScaleType string

const MetaDisplayScaleTypeLinear MetaDisplayScaleType = "linear"
const MetaDisplayScaleTypeLogarithmic MetaDisplayScaleType = "logarithmic"
const MetaDisplayScaleTypePower MetaDisplayScaleType = "power"
const MetaDisplayScaleTypeSquareroot MetaDisplayScaleType = "squareroot"

type MetaEmergencyMethodElem interface{}

type MetaProperties map[string]interface{}

type MetaWarnMethodElem interface{}

// A zone used to define the display and alarm state when the value is in between
// bottom and top.
type MetaZonesElem struct {
	// The lowest number in this zone
	Lower *float64 `json:"lower,omitempty" yaml:"lower,omitempty" mapstructure:"lower,omitempty"`

	// The message to display for the alarm.
	Message string `json:"message,omitempty" yaml:"message,omitempty" mapstructure:"message,omitempty"`

	// State corresponds to the JSON schema field "state".
	State AlarmState `json:"state" yaml:"state" mapstructure:"state"`

	// The highest value in this zone
	Upper *float64 `json:"upper,omitempty" yaml:"upper,omitempty" mapstructure:"upper,omitempty"`
}

// Maritime Mobile Service Identity (MMSI). Has to be 9 digits. See
// http://en.wikipedia.org/wiki/Maritime_Mobile_Service_Identity for information.
type Mmsi string

// Data should be of type NULL.
type NullValue struct {
	// Attr corresponds to the JSON schema field "_attr".
	Attr *Attr `json:"_attr,omitempty" yaml:"_attr,omitempty" mapstructure:"_attr,omitempty"`

	// Meta corresponds to the JSON schema field "meta".
	Meta *Meta `json:"meta,omitempty" yaml:"meta,omitempty" mapstructure:"meta,omitempty"`

	// Source corresponds to the JSON schema field "source".
	Source *Source `json:"source,omitempty" yaml:"source,omitempty" mapstructure:"source,omitempty"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp *Timestamp `json:"timestamp,omitempty" yaml:"timestamp,omitempty" mapstructure:"timestamp,omitempty"`

	// Value corresponds to the JSON schema field "value".
	Value interface{} `json:"value,omitempty" yaml:"value,omitempty" mapstructure:"value,omitempty"`
}

// Data should be of type number.
type NumberValue map[string]interface{}

// The position in 3 dimensions
type Position map[string]interface{}

// Resources to aid in navigation and operation of the vessel
type ResourcesJson struct {
	// A holder for charts, each named with their chart code
	Charts ResourcesJsonCharts `json:"charts,omitempty" yaml:"charts,omitempty" mapstructure:"charts,omitempty"`

	// A holder for notes about regions, each named with a UUID. Notes might include
	// navigation or cruising info, images, or anything
	Notes ResourcesJsonNotes `json:"notes,omitempty" yaml:"notes,omitempty" mapstructure:"notes,omitempty"`

	// A holder for regions, each named with UUID
	Regions ResourcesJsonRegions `json:"regions,omitempty" yaml:"regions,omitempty" mapstructure:"regions,omitempty"`

	// A holder for routes, each named with a UUID
	Routes ResourcesJsonRoutes `json:"routes,omitempty" yaml:"routes,omitempty" mapstructure:"routes,omitempty"`

	// A holder for waypoints, each named with a UUID
	Waypoints ResourcesJsonWaypoints `json:"waypoints,omitempty" yaml:"waypoints,omitempty" mapstructure:"waypoints,omitempty"`
}

// A holder for charts, each named with their chart code
type ResourcesJsonCharts map[string]interface{}

// A holder for notes about regions, each named with a UUID. Notes might include
// navigation or cruising info, images, or anything
type ResourcesJsonNotes map[string]interface{}

// A holder for regions, each named with UUID
type ResourcesJsonRegions map[string]interface{}

// A holder for routes, each named with a UUID
type ResourcesJsonRoutes map[string]interface{}

// A holder for waypoints, each named with a UUID
type ResourcesJsonWaypoints map[string]interface{}

// Maritime Mobile Service Identity (MMSI) for . Has to be 9 digits. See
// http://en.wikipedia.org/wiki/Maritime_Mobile_Service_Identity for information.
type SarMmsi string

// Root schema of Signal K. Contains the list of vessels plus a reference to the
// local boat (also contained in the vessels list).
type SignalkJson struct {
	// A wrapper object for aircraft, primarily intended for SAR aircraft in relation
	// to marine search and rescue. For clarity about seaplanes etc, if it CAN fly,
	// its an aircraft.
	Aircraft SignalkJsonAircraft `json:"aircraft,omitempty" yaml:"aircraft,omitempty" mapstructure:"aircraft,omitempty"`

	// A wrapper object for Aids to Navigation (aton's)
	Aton SignalkJsonAton `json:"aton,omitempty" yaml:"aton,omitempty" mapstructure:"aton,omitempty"`

	// Resources to aid in navigation and operation of the vessel including waypoints,
	// routes, notes, etc.
	Resources *ResourcesJson `json:"resources,omitempty" yaml:"resources,omitempty" mapstructure:"resources,omitempty"`

	// A wrapper object for Search And Rescue (SAR) MMSI's usied in transponders. MOB,
	// EPIRBS etc
	Sar SignalkJsonSar `json:"sar,omitempty" yaml:"sar,omitempty" mapstructure:"sar,omitempty"`

	// This holds the context (prefix + UUID, MMSI or URL in dot notation) of the
	// server's self object.
	Self string `json:"self" yaml:"self" mapstructure:"self"`

	// Metadata about the data sources; physical interface, address, protocol, etc.
	Sources SourcesJson `json:"sources,omitempty" yaml:"sources,omitempty" mapstructure:"sources,omitempty"`

	// Version of the schema and APIs that this data is using in Canonical format i.e.
	// V1.5.0.
	Version Version `json:"version" yaml:"version" mapstructure:"version"`

	// A wrapper object for vessel objects, each describing vessels in range,
	// including this vessel.
	Vessels SignalkJsonVessels `json:"vessels,omitempty" yaml:"vessels,omitempty" mapstructure:"vessels,omitempty"`
}

// A wrapper object for aircraft, primarily intended for SAR aircraft in relation
// to marine search and rescue. For clarity about seaplanes etc, if it CAN fly, its
// an aircraft.
type SignalkJsonAircraft map[string]interface{}

// A wrapper object for Aids to Navigation (aton's)
type SignalkJsonAton map[string]interface{}

// A wrapper object for Search And Rescue (SAR) MMSI's usied in transponders. MOB,
// EPIRBS etc
type SignalkJsonSar map[string]interface{}

// A wrapper object for vessel objects, each describing vessels in range, including
// this vessel.
type SignalkJsonVessels map[string]interface{}

// Source of data in delta format, a record of where the data was received from. An
// object containing at least the properties defined in 'properties', but can
// contain anything beyond that.
type Source struct {
	// AIS Message Type
	AisType *float64 `json:"aisType,omitempty" yaml:"aisType,omitempty" mapstructure:"aisType,omitempty"`

	// NMEA2000 can name of the source device
	CanName *string `json:"canName,omitempty" yaml:"canName,omitempty" mapstructure:"canName,omitempty"`

	// NMEA2000 instance value of the source message
	Instance *string `json:"instance,omitempty" yaml:"instance,omitempty" mapstructure:"instance,omitempty"`

	// A label to identify the source bus, eg serial-COM1, eth-local,etc . Can be
	// anything but should follow a predicatable format
	Label string `json:"label" yaml:"label" mapstructure:"label"`

	// NMEA2000 pgn of the source message
	Pgn *float64 `json:"pgn,omitempty" yaml:"pgn,omitempty" mapstructure:"pgn,omitempty"`

	// Sentence type of the source NMEA0183 sentence,
	// $GP[RMC],092750.000,A,5321.6802,N,00630.3372,W,0.02,31.66,280511,,,A*43
	Sentence *string `json:"sentence,omitempty" yaml:"sentence,omitempty" mapstructure:"sentence,omitempty"`

	// NMEA2000 src value or any similar value for encapsulating the original source
	// of the data
	Src *string `json:"src,omitempty" yaml:"src,omitempty" mapstructure:"src,omitempty"`

	// Talker id of the source NMEA0183 sentence,
	// $[GP]RMC,092750.000,A,5321.6802,N,00630.3372,W,0.02,31.66,280511,,,A*43
	Talker *string `json:"talker,omitempty" yaml:"talker,omitempty" mapstructure:"talker,omitempty"`

	// A human name to identify the type. NMEA0183, NMEA2000, signalk
	Type string `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`
}

// Reference to the source under /sources. A dot spearated path to the data. eg
// [type].[bus].[device]
type SourceRef string

// Metadata about the sources, eg. buses and connected sensors
type SourcesJson map[string]interface{}

// Data should be of type number.
type StringValue map[string]interface{}

// RFC 3339 (UTC only without local offset) string representing date and time.
type Timestamp string

// Allowed units of physical quantities. Units should be (derived) SI units where
// possible.
type Units string

// UnmarshalJSON implements json.Unmarshaler.
func (j *HelloJsonRolesElem) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_HelloJsonRolesElem {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_HelloJsonRolesElem, v)
	}
	*j = HelloJsonRolesElem(v)
	return nil
}

type ValuesDatetimeValue struct {
	// Pgn corresponds to the JSON schema field "pgn".
	Pgn *float64 `json:"pgn,omitempty" yaml:"pgn,omitempty" mapstructure:"pgn,omitempty"`

	// Sentence corresponds to the JSON schema field "sentence".
	Sentence *string `json:"sentence,omitempty" yaml:"sentence,omitempty" mapstructure:"sentence,omitempty"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp *Timestamp `json:"timestamp,omitempty" yaml:"timestamp,omitempty" mapstructure:"timestamp,omitempty"`

	// Value corresponds to the JSON schema field "value".
	Value *string `json:"value,omitempty" yaml:"value,omitempty" mapstructure:"value,omitempty"`
}

type ValuesNumberValue struct {
	// Pgn corresponds to the JSON schema field "pgn".
	Pgn *float64 `json:"pgn,omitempty" yaml:"pgn,omitempty" mapstructure:"pgn,omitempty"`

	// Sentence corresponds to the JSON schema field "sentence".
	Sentence *string `json:"sentence,omitempty" yaml:"sentence,omitempty" mapstructure:"sentence,omitempty"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp *Timestamp `json:"timestamp,omitempty" yaml:"timestamp,omitempty" mapstructure:"timestamp,omitempty"`

	// Value corresponds to the JSON schema field "value".
	Value *float64 `json:"value,omitempty" yaml:"value,omitempty" mapstructure:"value,omitempty"`
}

type ValuesStringValue struct {
	// Pgn corresponds to the JSON schema field "pgn".
	Pgn *float64 `json:"pgn,omitempty" yaml:"pgn,omitempty" mapstructure:"pgn,omitempty"`

	// Sentence corresponds to the JSON schema field "sentence".
	Sentence *string `json:"sentence,omitempty" yaml:"sentence,omitempty" mapstructure:"sentence,omitempty"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp *Timestamp `json:"timestamp,omitempty" yaml:"timestamp,omitempty" mapstructure:"timestamp,omitempty"`

	// Value corresponds to the JSON schema field "value".
	Value *string `json:"value,omitempty" yaml:"value,omitempty" mapstructure:"value,omitempty"`
}

// Version of the Signal K schema/APIs used by the root object.
type Version string

// A waypoint, an object with a signal k position object, and GeoJSON Feature
// object (see geojson.org, and
// https://github.com/fge/sample-json-schemas/tree/master/geojson)
type Waypoint struct {
	// A Geo JSON feature object
	Feature interface{} `json:"feature,omitempty" yaml:"feature,omitempty" mapstructure:"feature,omitempty"`

	// Position corresponds to the JSON schema field "position".
	Position Position `json:"position,omitempty" yaml:"position,omitempty" mapstructure:"position,omitempty"`
}

var enumValues_AlarmMethodEnum = []interface{}{
	"visual",
	"sound",
}

// A location of a resource, potentially relative. For hierarchical schemes (like
// http), applications must resolve relative URIs (e.g. './v1/api/').
// Implementations should support the following schemes: http:, https:, mailto:,
// tel:, and ws:.
type Url string

// UnmarshalJSON implements json.Unmarshaler.
func (j *AlarmMethodEnum) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_AlarmMethodEnum {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_AlarmMethodEnum, v)
	}
	*j = AlarmMethodEnum(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SignalkJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["self"]; !ok || v == nil {
		return fmt.Errorf("field self in SignalkJson: required")
	}
	if v, ok := raw["version"]; !ok || v == nil {
		return fmt.Errorf("field version in SignalkJson: required")
	}
	type Plain SignalkJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = SignalkJson(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Attr) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	type Plain Attr
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["_group"]; !ok || v == nil {
		plain.Group = "self"
	}
	if v, ok := raw["_mode"]; !ok || v == nil {
		plain.Mode = 644.0
	}
	if v, ok := raw["_owner"]; !ok || v == nil {
		plain.Owner = "self"
	}
	*j = Attr(plain)
	return nil
}

var enumValues_HelloJsonRolesElem = []interface{}{
	"master",
	"main",
	"aux",
	"slave",
}

// A unique Signal K flavoured maritime resource identifier (MRN). A MRN is a form
// of URN, following a specific format: urn:mrn:<issueing authority>:<id
// type>:<id>. In case of a Signal K uuid, that looks like this:
// urn:mrn:signalk:uuid:<uuid>, where Signal K is the issuing authority and UUID
// (v4) the ID type.
type Uuid string

// UnmarshalJSON implements json.Unmarshaler.
func (j *AlarmState) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_AlarmState {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_AlarmState, v)
	}
	*j = AlarmState(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CommonValueFields) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["$source"]; !ok || v == nil {
		return fmt.Errorf("field $source in CommonValueFields: required")
	}
	if v, ok := raw["timestamp"]; !ok || v == nil {
		return fmt.Errorf("field timestamp in CommonValueFields: required")
	}
	type Plain CommonValueFields
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CommonValueFields(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *NullValue) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	type Plain NullValue
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Value != nil {
		return fmt.Errorf("field %s: must be null", "value")
	}
	*j = NullValue(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Meta) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["description"]; !ok || v == nil {
		return fmt.Errorf("field description in Meta: required")
	}
	type Plain Meta
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["alarmMethod"]; !ok || v == nil {
		plain.AlarmMethod = []MetaAlarmMethodElem{
			"visual",
			"sound",
		}
	}
	if v, ok := raw["alertMethod"]; !ok || v == nil {
		plain.AlertMethod = []MetaAlertMethodElem{
			"visual",
		}
	}
	if v, ok := raw["emergencyMethod"]; !ok || v == nil {
		plain.EmergencyMethod = []MetaEmergencyMethodElem{
			"visual",
			"sound",
		}
	}
	if v, ok := raw["warnMethod"]; !ok || v == nil {
		plain.WarnMethod = []MetaWarnMethodElem{
			"visual",
		}
	}
	*j = Meta(plain)
	return nil
}

var enumValues_AlarmState = []interface{}{
	"nominal",
	"normal",
	"alert",
	"warn",
	"alarm",
	"emergency",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *HelloJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["roles"]; !ok || v == nil {
		return fmt.Errorf("field roles in HelloJson: required")
	}
	if v, ok := raw["version"]; !ok || v == nil {
		return fmt.Errorf("field version in HelloJson: required")
	}
	type Plain HelloJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if len(plain.Roles) < 2 {
		return fmt.Errorf("field %s length: must be >= %d", "roles", 2)
	}
	if len(plain.Roles) > 2 {
		return fmt.Errorf("field %s length: must be <= %d", "roles", 2)
	}
	*j = HelloJson(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *MetaZonesElem) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["state"]; !ok || v == nil {
		return fmt.Errorf("field state in MetaZonesElem: required")
	}
	type Plain MetaZonesElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["message"]; !ok || v == nil {
		plain.Message = "Warning"
	}
	*j = MetaZonesElem(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeltaJsonUpdatesElemMetaElem) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["path"]; !ok || v == nil {
		return fmt.Errorf("field path in DeltaJsonUpdatesElemMetaElem: required")
	}
	if v, ok := raw["value"]; !ok || v == nil {
		return fmt.Errorf("field value in DeltaJsonUpdatesElemMetaElem: required")
	}
	type Plain DeltaJsonUpdatesElemMetaElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeltaJsonUpdatesElemMetaElem(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *MetaDisplayScaleType) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_MetaDisplayScaleType {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_MetaDisplayScaleType, v)
	}
	*j = MetaDisplayScaleType(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeltaJsonUpdatesElemValuesElem) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["path"]; !ok || v == nil {
		return fmt.Errorf("field path in DeltaJsonUpdatesElemValuesElem: required")
	}
	if v, ok := raw["value"]; !ok || v == nil {
		return fmt.Errorf("field value in DeltaJsonUpdatesElemValuesElem: required")
	}
	type Plain DeltaJsonUpdatesElemValuesElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeltaJsonUpdatesElemValuesElem(plain)
	return nil
}

var enumValues_MetaDisplayScaleType = []interface{}{
	"linear",
	"logarithmic",
	"squareroot",
	"power",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Source) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["label"]; !ok || v == nil {
		return fmt.Errorf("field label in Source: required")
	}
	type Plain Source
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["type"]; !ok || v == nil {
		plain.Type = "NMEA2000"
	}
	*j = Source(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeltaJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["updates"]; !ok || v == nil {
		return fmt.Errorf("field updates in DeltaJson: required")
	}
	type Plain DeltaJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeltaJson(plain)
	return nil
}
