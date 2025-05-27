package chat

import (
	"context"
	"encoding/json"
	"log"
)

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte, 100), // Буферизованный канал
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	log.Println("Hub started running")
	for {
		select {
		case client := <-h.Register:
			log.Printf("[HUB] Registering new client: UserID=%d, Username=%s", client.UserID, client.Username)
			h.Clients[client] = true

			messages, err := client.ChatUC.GetRecentMessages(context.Background(), 50)
			if err != nil {
				log.Printf("[HUB] Error getting messages: %v", err)
				continue
			}
			log.Printf("[HUB] Sending %d historical messages to client %d", len(messages), client.UserID)

			for _, msg := range messages {
				jsonMsg, err := json.Marshal(msg)
				if err != nil {
					log.Printf("[HUB] Error marshaling message: %v", err)
					continue
				}
				select {
				case client.Send <- jsonMsg:
					log.Printf("[HUB] Historical message sent to %d", client.UserID)
				default:
					log.Printf("[HUB] Client %d send channel blocked, closing", client.UserID)
					close(client.Send)
					delete(h.Clients, client)
				}
			}

		case client := <-h.Unregister:
			log.Printf("[HUB] Unregistering client: UserID=%d", client.UserID)
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}

		case message := <-h.Broadcast:
			log.Printf("[HUB] Broadcasting message to %d clients: %s", len(h.Clients), string(message))
			if len(message) == 0 {
				log.Println("[HUB] Warning: empty message received")
				continue
			}

			for client := range h.Clients {
				select {
				case client.Send <- message:
					log.Printf("[HUB] Message sent to client %d", client.UserID)
				default:
					log.Printf("[HUB] Client %d channel blocked, disconnecting", client.UserID)
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
