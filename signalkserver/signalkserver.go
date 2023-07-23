package signalkserver

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/wdantuma/signalk-server-go/converter"
	"github.com/wdantuma/signalk-server-go/resources/charts"
	"github.com/wdantuma/signalk-server-go/source"
	"github.com/wdantuma/signalk-server-go/store"
	"github.com/wdantuma/signalk-server-go/stream"
	"github.com/wdantuma/signalk-server-go/vessel"
)

var Version = "0.0.1" // overwritten with VERSION DEF during build

const (
	SERVER_NAME string = "signalk-server-go"
)

type signalkServer struct {
	name      string
	version   string
	self      string
	debug     bool
	store     store.Store
	sourcehub *source.Sourcehub
	converter converter.CanToSignalk
}

func NewSignalkServer() *signalkServer {
	self := fmt.Sprintf("vessels.urn:mrn:signalk:uuid:%s", uuid.New().String())
	return &signalkServer{name: SERVER_NAME, version: Version, self: self, sourcehub: source.NewSourceHub()}
}

func (s *signalkServer) GetName() string {
	return s.name
}

func (s *signalkServer) GetVersion() string {
	return s.version
}

func (s *signalkServer) GetSelf() string {
	return s.self
}

func (s *signalkServer) GetDebug() bool {
	return s.debug
}

func (s *signalkServer) SetDebug(debug bool) {
	s.debug = debug
}

func (s *signalkServer) GetStore() store.Store {
	return s.store
}

func (s *signalkServer) SetMMSI(mmsi string) {
	s.self = fmt.Sprintf("vessels.urn:mrn:imo:mmsi:%s", mmsi)
}

func (server *signalkServer) AddSource(source source.CanSource) {
	server.sourcehub.AddSource(source)
}

func (server *signalkServer) hello(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	method := "http"
	wsmethod := "ws"
	if req.TLS != nil {
		method = "https"
		wsmethod = "wss"

	}
	fmt.Fprintf(w, `
	{
		"endpoints": {
			"v1": {
				"version": "2.0.0",
				"signalk-http": "%s://%s/signalk/v1/api/",
				"signalk-ws": "%s://%s/signalk/v1/stream"
			}
		},
		"server": {
			"id": "signalk-server-go",
			"version": "%s"
		}
	}
`, method, req.Host, wsmethod, req.Host, server.GetVersion())
}

func (server *signalkServer) loginStatus(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `
	{
		"status": "notLoggedIn",
		"readOnlyAccess": true,
		"authenticationRequired": true,
		"allowNewUserRegistration": true,
		"allowDeviceAccessRequests": true,
		"securityWasEnabled": false
	}
`)
}

func (server *signalkServer) SetupServer(ctx context.Context, hostname string, router *mux.Router) *mux.Router {
	var err error
	if router == nil {
		router = mux.NewRouter()
	}

	signalk := router.PathPrefix("/signalk").Subrouter()
	signalk.HandleFunc("", server.hello)
	streamHandler := stream.NewStreamHandler(server)
	vesselHandler := vessel.NewVesselHandler(server)
	chartsHandler := charts.NewChartsHandler()
	signalk.PathPrefix("/v1/stream").Handler(streamHandler)
	signalk.PathPrefix("/v1/api/vessels").Handler(vesselHandler)
	signalk.PathPrefix("/v1/api/resources/charts").Handler(chartsHandler)

	signalk.HandleFunc("/v1/api/snapshot", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})

	router.HandleFunc("/skServer/loginStatus", server.loginStatus)

	server.converter, err = converter.NewCanToSignalk(server)
	if err != nil {
		log.Fatal(err)
	}

	canSource := server.sourcehub.Start()
	converted := server.converter.Convert(canSource)
	valueStore := store.NewMemoryStore()
	server.store = valueStore
	stored := valueStore.Store(converted)

	go func() {
		for delta := range stored {
			streamHandler.BroadcastDelta <- delta
		}
	}()

	return router
}
