package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type User struct {
	Name string
	Conn *websocket.Conn
}

// Represent incoming messages from users to the hub.
type InMessage struct {
	Content string `json:"content"`
}

// Represent outgoing messages broadcasting from the hub.
type OutMessage struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

type ChatHub struct {
	users    map[string]User // map username -> User
	mu       sync.Mutex      // protect users map
	messages chan OutMessage // broadcast channel from one to other users
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		users:    make(map[string]User),
		messages: make(chan OutMessage, 1000),
		mu:       sync.Mutex{},
	}
}

func (h *ChatHub) AddUser(u User) {
	h.mu.Lock()
	h.users[u.Name] = u
	h.mu.Unlock()
}

func (h *ChatHub) DelUser(name string) {
	h.mu.Lock()
	delete(h.users, name)
	h.mu.Unlock()
}

// Handle route /ws?name={}
func (hub *ChatHub) handleWS(c *gin.Context) {
	// Get name from query string.
	name := c.Query("name")
	if name == "" {
		c.String(http.StatusBadRequest, "missing query parameter `name`")
		return
	}

	// Upgrade http connection to WS.
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	user := User{name, conn}
	hub.AddUser(user)

	// Reading incoming messages until connection is closed.
	// Send messages to an internal broadcast channel.
	for {
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			hub.DelUser(user.Name)
			conn.Close()
			return
		}

		var inMsg InMessage
		if err := json.Unmarshal(bytes, &inMsg); err != nil {
			log.Println("Decode incoming JSON failed:", err)
			return
		}
		outMsg := OutMessage{
			Content:  inMsg.Content,
			Username: name,
		}
		hub.messages <- outMsg
	}
}

// Broadcast incoming messages from the internal broadcast channel
func (h *ChatHub) broadcast() {
	for msg := range h.messages {
		bytes, err := json.Marshal(&msg)
		if err != nil {
			log.Println("Encode outgoing message failed:", err)
			continue
		}

		// Copy users to avoid race condition.
		h.mu.Lock()
		curUsers := make([]User, 0, len(h.users))
		for _, u := range h.users {
			curUsers = append(curUsers, u)
		}
		h.mu.Unlock()

		for _, user := range curUsers {
			// Do not echo back to message owner.
			if user.Name == msg.Username {
				continue
			}
			// Broadcast concurrently.
			go func() {
				if err := user.Conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
					h.DelUser(user.Name)
					user.Conn.Close()
				}
			}()
		}
	}
}

func main() {
	router := gin.Default()
	hub := NewChatHub()

	// Register routes with handlers.
	router.GET("/ws", hub.handleWS)

	router.GET("health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Start broadcasting.
	go hub.broadcast()

	router.Run(":5001")
}
