package state

const (
	SERVER_NAME string = "signalk-server-go"
	TIME_FORMAT string = "2006-01-02T15:04:05.000Z"
	SELF        string = "vessels.urn:mrn:imo:mmsi:244810236" //244810236
	VERSION     string = "0.0.1"
)

type ServerState interface {
	GetName() string
	GetVersion() string
	GetSelf() string
}
