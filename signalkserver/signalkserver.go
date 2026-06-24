package signalkserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/wdantuma/signalk-server-go/resources/charts"
	"github.com/wdantuma/signalk-server-go/source"
	"github.com/wdantuma/signalk-server-go/source/base"
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

	chartsPath string
}

func NewSignalkServer() *signalkServer {
	self := fmt.Sprintf("vessels.urn:mrn:signalk:uuid:%s", uuid.New().String())
	server := &signalkServer{name: SERVER_NAME, version: Version, self: self}

	server.sourcehub = source.NewSourceHub()
	return server
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

func (s *signalkServer) SetChartsPath(chartsPath string) {
	s.chartsPath = chartsPath
}

func (s *signalkServer) GetStore() store.Store {
	return s.store
}

func (s *signalkServer) SetMMSI(mmsi string) {
	s.self = fmt.Sprintf("vessels.urn:mrn:imo:mmsi:%s", mmsi)
}

func (server *signalkServer) AddSource(source base.DeltaSource) {
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

func (server *signalkServer) features(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `
	{
		"apis": [
		   "resources"
		 ],
		 "plugins":[]
	}
`)
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
	if router == nil {
		router = mux.NewRouter()
	}

	signalk := router.PathPrefix("/signalk").Subrouter()
	streamHandler := stream.NewStreamHandler(server)
	vesselHandler := vessel.NewVesselHandler(server)
	chartsHandler := charts.NewChartsHandler(server.chartsPath)
	signalk.PathPrefix("/v1/stream").Handler(streamHandler)
	signalk.HandleFunc("/v2/features", server.features)
	signalk.PathPrefix("/v1/api/vessels").Handler(vesselHandler)
	signalk.PathPrefix("/v2/api/resources/charts").Handler(chartsHandler)
	signalk.PathPrefix("/v1/api/resources/charts").Handler(chartsHandler)
	signalk.HandleFunc("/v1/api/snapshot", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})
	signalk.HandleFunc("/", server.hello)
	signalk.HandleFunc("", server.hello)

	router.HandleFunc("/skServer/loginStatus", server.loginStatus)

	hub := server.sourcehub.Start()
	valueStore := store.NewMemoryStore()
	server.store = valueStore
	stored := valueStore.Store(hub)

	go func() {
		for delta := range stored {
			streamHandler.BroadcastDelta <- delta
		}
	}()

	return router
}
