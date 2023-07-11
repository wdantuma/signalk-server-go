package stream

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wdantuma/signalk-server-go/ref"
	"github.com/wdantuma/signalk-server-go/signalk"
	"github.com/wdantuma/signalk-server-go/signalk/filter"
	"github.com/wdantuma/signalk-server-go/signalk/format"

	//	"github.com/wdantuma/signalk-server-go/signalkserver"
	"github.com/wdantuma/signalk-server-go/signalkserver/state"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

type streamHandler struct {
	state state.ServerState
}

func NewStreamHandler(s state.ServerState) *streamHandler {
	return &streamHandler{state: s}
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type client struct {
	filter *filter.Filter
	hub    *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send      chan []byte
	sendDelta chan signalk.DeltaJson
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

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		if len(message) > 3 {
			subscribeMessage := signalk.SubscribeJson{}
			err := subscribeMessage.UnmarshalJSON(message)
			if err == nil {
				c.filter.UpdateSubscription(subscribeMessage)
			}
		}

		//c.hub.Broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func (s *streamHandler) ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	contextFilter := filter.NewFilter(s.state.GetSelf())
	contextFilter.Subscribe = filter.ParseSubscribe(r.URL.Query().Get("subscribe"))
	client := &client{hub: hub, filter: contextFilter, conn: conn, send: make(chan []byte, 256), sendDelta: make(chan signalk.DeltaJson)}
	client.hub.register <- client

	client.send <- s.helloMessage()
	format.Json(contextFilter.Filter(client.sendDelta), client.send)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
