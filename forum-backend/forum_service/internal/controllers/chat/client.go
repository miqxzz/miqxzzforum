package chat

import (
	"context"
	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/internal/usecase"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	Hub             *Hub
	Conn            *websocket.Conn
	Send            chan []byte
	UserID          int
	Username        string
	IsAuthenticated bool
	ChatUC          usecase.ChatUsecase
}

func (c *Client) ReadPump() {
	log.Printf("[CLIENT %d] Starting read pump", c.UserID)
	defer func() {
		log.Printf("[CLIENT %d] Closing read pump", c.UserID)
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, rawMessage, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[CLIENT %d] Read error: %v", c.UserID, err)
			}
			break
		}
		log.Printf("[CLIENT %d] Received raw message: %s", c.UserID, string(rawMessage))

		if err2 := c.handleIncomingMessage(rawMessage); err2 != nil {
			log.Printf("[CLIENT %d] Message handling error: %v", c.UserID, err2)
		}
	}
}

func (c *Client) handleIncomingMessage(rawMessage []byte) error {
	if len(rawMessage) == 0 {
		log.Printf("[CLIENT %d] Empty message received", c.UserID)
		return nil
	}

	var msg entity.ChatMessage
	if err := json.Unmarshal(rawMessage, &msg); err != nil {
		log.Printf("[CLIENT %d] Failed to unmarshal message, creating default: %v", c.UserID, err)
		msg = entity.ChatMessage{
			UserID:    c.UserID,
			Username:  c.Username,
			Content:   string(rawMessage),
			Timestamp: time.Now(),
		}
	}

	if msg.Content == "" {
		log.Printf("[CLIENT %d] Empty content in message", c.UserID)
		return nil
	}
	log.Println(c.IsAuthenticated)
	if c.IsAuthenticated {
		log.Printf("[CLIENT %d] Saving message to DB: %s", c.UserID, msg.Content)
		if err := c.ChatUC.HandleMessage(context.Background(), msg.UserID, msg.Username, msg.Content); err != nil {
			log.Printf("[CLIENT %d] DB save error: %v", c.UserID, err)
		}
	}

	outMsg := map[string]interface{}{
		"userID":    msg.UserID,
		"username":  msg.Username,
		"content":   msg.Content,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	jsonMsg, err := json.Marshal(outMsg)
	if err != nil {
		log.Printf("[CLIENT %d] Marshal error: %v", c.UserID, err)
		return err
	}

	log.Printf("[CLIENT %d] Broadcasting message: %s", c.UserID, string(jsonMsg))
	c.Hub.Broadcast <- jsonMsg
	return nil
}

func (c *Client) WritePump() {
	log.Printf("[CLIENT %d] Starting write pump", c.UserID)
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		log.Printf("[CLIENT %d] Closing write pump", c.UserID)
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				log.Printf("[CLIENT %d] Send channel closed, sending close message", c.UserID)
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			log.Printf("[CLIENT %d] Preparing to write message: %s", c.UserID, string(message))
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("[CLIENT %d] NextWriter error: %v", c.UserID, err)
				return
			}
			if _, err := w.Write(message); err != nil {
				log.Printf("[CLIENT %d] Write error: %v", c.UserID, err)
				return
			}

			if err := w.Close(); err != nil {
				log.Printf("[CLIENT %d] Writer close error: %v", c.UserID, err)
				return
			}
			log.Printf("[CLIENT %d] Message successfully sent", c.UserID)

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[CLIENT %d] Ping error: %v", c.UserID, err)
				return
			}
			log.Printf("[CLIENT %d] Ping sent", c.UserID)
		}
	}
}
