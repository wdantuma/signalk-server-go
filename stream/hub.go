package stream

type Hub struct {
	// Registered clients.
	clients map[*client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	register chan *client

	// Unregister requests from clients.
	unregister chan *client
}

func NewHub() *Hub {
	hub := &Hub{
		Broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
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
					close(client.send)
				}
			case message := <-h.Broadcast:
				for client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}()
}
