package game

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	GameID   string
	PlayerID string
	Color    string
	Send     chan []byte
}

type Hub struct {
	rooms    map[string]map[*Client]bool
	mu       sync.RWMutex
	register chan *Client
	unregister chan *Client
	broadcast chan *Message
}

type Message struct {
	GameID string
	Data   []byte
}

func NewHub() *Hub {
	return &Hub{
		rooms:     make(map[string]map[*Client]bool),
		register:  make(chan *Client),
		unregister: make(chan *Client),
		broadcast: make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.rooms[client.GameID]; !ok {
				h.rooms[client.GameID] = make(map[*Client]bool)
			}
			h.rooms[client.GameID][client] = true
			h.mu.Unlock()

			log.Printf("Client %s joined game %s", client.PlayerID, client.GameID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.GameID]; ok {
				delete(clients, client)
				if len(clients) == 0 {
					delete(h.rooms, client.GameID)
				}
			}
			h.mu.Unlock()

			log.Printf("Client %s left game %s", client.PlayerID, client.GameID)

		case msg := <-h.broadcast:
			h.mu.RLock()
			if clients, ok := h.rooms[msg.GameID]; ok {
				for client := range clients {
					select {
					case client.Send <- msg.Data:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(gameID string, data []byte) {
	h.broadcast <- &Message{GameID: gameID, Data: data}
}
