package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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

func stream(ws *websocket.Conn) {
	defer func() {
		ws.Close()
	}()
	for {
		// TODO generate signal-k messages
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	go stream(ws)

}

func main() {
	route := mux.NewRouter()
	fs := http.FileServer(http.Dir("./static"))
	route.HandleFunc("/signalk/v1/stream", serveWs)
	route.HandleFunc("/signalk", signalk)
	route.PathPrefix("/@signalk").Handler(fs)

	route.Handle("/", http.RedirectHandler("/@signalk/freeboard-sk", http.StatusSeeOther))

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", route)
	if err != nil {
		log.Fatal(err)
	}
}
