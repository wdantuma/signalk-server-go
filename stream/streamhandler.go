package stream

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/wdantuma/signalk-server-go/ref"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/signalk/filter"
	"github.com/wdantuma/signalk-server-go/signalk/format"
	"github.com/wdantuma/signalk-server-go/signalkserver/state"
)

type streamHandler struct {
	state          state.ServerState
	BroadcastDelta chan signalk.DeltaJson
	hub            *hub
}

func NewStreamHandler(s state.ServerState) *streamHandler {
	hub := NewHub()
	return &streamHandler{state: s, hub: hub, BroadcastDelta: hub.BroadcastDelta}
}

func (s *streamHandler) helloMessage() []byte {
	hello := signalk.HelloJson{}
	hello.Name = ref.String(s.state.GetName())
	hello.Version = (signalk.Version)(s.state.GetVersion())
	hello.Timestamp = ref.UTCTimeStamp(time.Now())
	hello.Self = ref.String(s.state.GetSelf())
	hello.Roles = append(hello.Roles, "master")
	hello.Roles = append(hello.Roles, "main")
	helloBytes, _ := json.Marshal(hello)
	return helloBytes
}

// serveWs handles websocket requests from the peer.
func (s *streamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	contextFilter := filter.NewFilter(s.state.GetSelf())
	contextFilter.Subscribe = filter.ParseSubscribe(r.URL.Query().Get("subscribe"))
	client := &client{hub: s.hub, filter: contextFilter, conn: conn, send: make(chan []byte, 1024), sendDelta: make(chan signalk.DeltaJson, 10)}
	format.Json(contextFilter.Filter(client.sendDelta), client.send)
	time.Sleep(1 * time.Second)
	client.hub.register <- client

	client.send <- s.helloMessage()

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
