package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// Client struct to hold the WebSocket connection and a channel for messages
type Client struct {
	Conn *websocket.Conn
	send chan []byte
	username string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	// continuously listens for events on the register, unregister and broadcase
	for {
		select {
		case client := <-h.register:
			print(client)
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				// sends message to the client.send channel
				case client.send <- message:
				default: //channel is blocked
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func HandleConnections(h *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := &Client{
		Conn: conn,
		send: make(chan []byte),
	}
	h.register <- client
	go client.readMessages(h)
	go client.writeMessage()
}

func (c *Client) readMessages(h *Hub) {
	defer func() {
		h.unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		h.broadcast <- message
	}
}

func (c *Client) writeMessage() {
	defer c.Conn.Close()
	for message := range c.send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
	}
}
