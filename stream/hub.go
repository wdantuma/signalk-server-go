package stream

import (
	"github.com/wdantuma/signalk-server-go/signalk"
)

type Hub struct {
	// Registered clients.
	clients map[*client]bool

	// Inbound messages from the clients.
	BroadcastDelta chan signalk.DeltaJson

	// Register requests from the clients.
	register chan *client

	// Unregister requests from clients.
	unregister chan *client
}

func NewHub() *Hub {
	hub := &Hub{
		BroadcastDelta: make(chan signalk.DeltaJson),
		register:       make(chan *client),
		unregister:     make(chan *client),
		clients:        make(map[*client]bool),
	}
	hub.run()
	return hub
}

func (h *Hub) run() {
	go func() {
		for {
			select {
			case client := <-h.register:
				h.clients[client] = true
			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.sendDelta)
				}
			case message := <-h.BroadcastDelta:
				for client := range h.clients {
					select {
					case client.sendDelta <- message:
					default:
						close(client.sendDelta)
						delete(h.clients, client)
					}
				}
			}
		}
	}()
}
