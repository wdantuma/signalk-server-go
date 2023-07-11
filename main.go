package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/wdantuma/signalk-server-go/converter"
	"github.com/wdantuma/signalk-server-go/socketcan"
	"github.com/wdantuma/signalk-server-go/stream"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func signalk(w http.ResponseWriter, req *http.Request) {
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
			"version": "2.0.0"
		}
	}
`, method, req.Host, wsmethod, req.Host)
}

func main() {
	hub := stream.NewHub()

	route := mux.NewRouter()
	fs := http.FileServer(http.Dir("./static"))
	route.HandleFunc("/signalk/v1/stream", func(w http.ResponseWriter, r *http.Request) {
		stream.ServeWs(hub, w, r)
	})
	route.HandleFunc("/signalk", signalk)
	route.PathPrefix("/@signalk").Handler(fs)

	route.Handle("/", http.RedirectHandler("/@signalk/instrumentpanel", http.StatusSeeOther))

	// main loop
	source, err := socketcan.NewCanDumpSource("data/n2kdump.txt")
	if err != nil {
		log.Fatal(err)
	}
	converter, err := converter.NewCanToSignalk()
	if err != nil {
		log.Fatal(err)
	}

	sk := converter.Convert(source)

	go func() {
		for bytes := range sk {
			hub.BroadcastDelta <- bytes
		}
	}()

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", route)
	if err != nil {
		log.Fatal(err)
	}
}
