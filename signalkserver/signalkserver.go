package signalkserver

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/wdantuma/signalk-server-go/converter"
	"github.com/wdantuma/signalk-server-go/socketcan"
	"github.com/wdantuma/signalk-server-go/stream"
)

var Version = "0.0.1"

const (
	SERVER_NAME string = "signalk-server-go"
	TIME_FORMAT string = "2006-01-02T15:04:05.000Z"
)

type signalkServer struct {
	name    string
	version string
	self    string
	debug   bool
}

func NewSignalkServer() *signalkServer {
	self := fmt.Sprintf("vessels.urn:mrn:signalk:uuid:%s", uuid.New().String())
	return &signalkServer{name: SERVER_NAME, version: Version, self: self}
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

func (s *signalkServer) EnableDebug() {
	s.debug = true
}

func (s *signalkServer) SetMMSI(mmsi string) {
	s.self = fmt.Sprintf("vessels.urn:mrn:imo:mmsi:%s", mmsi)
}

func (server *signalkServer) Hello(w http.ResponseWriter, req *http.Request) {
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

func (server *signalkServer) SetupServer(ctx context.Context, hostname string, router *mux.Router) *mux.Router {
	if router == nil {
		router = mux.NewRouter()
	}

	hub := stream.NewHub()

	signalk := router.PathPrefix("/signalk").Subrouter()
	signalk.HandleFunc("", server.Hello)
	streamHandler := stream.NewStreamHandler(server)
	signalk.HandleFunc("/v1/stream", func(w http.ResponseWriter, r *http.Request) {
		streamHandler.ServeWs(hub, w, r)
	})

	// main loop
	source, err := socketcan.NewCanDumpSource("data/n2kdump.txt")
	if err != nil {
		log.Fatal(err)
	}
	converter, err := converter.NewCanToSignalk()
	if err != nil {
		log.Fatal(err)
	}

	sk := converter.Convert(server, source)

	go func() {
		for bytes := range sk {
			hub.BroadcastDelta <- bytes
		}
	}()

	return router
}
