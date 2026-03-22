package websocket

import (
	"code-mafia-backend/internal/database"
	"code-mafia-backend/internal/redis"
	"encoding/json"
	"log"
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	redis      *redis.Client
	repo       *database.Repository
	secretKey  string
	mu         sync.RWMutex
}

func NewHub(redisClient *redis.Client, repo *database.Repository, secretKey string) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		redis:      redisClient,
		repo:       repo,
		secretKey:  secretKey,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered: %s (ID: %s)", client.username, client.id)
			h.broadcastUserList()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client unregistered: %s (ID: %s)", client.username, client.id)
			}
			h.mu.Unlock()
			h.broadcastUserList()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) broadcastUserList() {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Deduplicate by username — multiple connections from same team count as one
	seen := make(map[string]bool)
	users := make([]map[string]string, 0)
	for client := range h.clients {
		if seen[client.username] {
			continue
		}
		seen[client.username] = true
		users = append(users, map[string]string{
			"userID":   client.id,
			"username": client.username,
		})
	}

	message := map[string]interface{}{
		"type":    "users",
		"payload": users,
	}

	h.BroadcastMessage(message)
}

// sendUserListTo sends the current user list to a specific client, excluding themselves.
func (h *Hub) sendUserListTo(c *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Deduplicate by username and exclude the requesting client's own team
	seen := make(map[string]bool)
	seen[c.username] = true // pre-exclude self
	users := make([]map[string]string, 0)
	for client := range h.clients {
		if seen[client.username] {
			continue
		}
		seen[client.username] = true
		users = append(users, map[string]string{
			"userID":   client.id,
			"username": client.username,
		})
	}
	c.sendMessage("users", users)
}

func (h *Hub) BroadcastMessage(message interface{}) {
	data, _ := json.Marshal(message)
	h.broadcast <- data
}

func (h *Hub) GetClientByID(id string) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.id == id {
			return client
		}
	}
	return nil
}

func (h *Hub) GetClientByUsername(username string) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.username == username {
			return client
		}
	}
	return nil
}
