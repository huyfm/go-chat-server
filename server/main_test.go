package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestChatHub(t *testing.T) {
	hub := NewChatHub()
	if len(hub.users) != 0 {
		t.Fatal("invalid initial number of users")
	}

	hub.AddUser(User{Name: "foo"})
	if len(hub.users) != 1 {
		t.Fatal("invalid number of users after adding one")
	}

	hub.DelUser("foo")
	if len(hub.users) != 0 {
		t.Fatal("invalid number of users after deleting one")
	}
}

func TestHandleWS(t *testing.T) {
	// Setup test server.
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	hub := NewChatHub()
	router.GET("/ws", hub.handleWS)
	go hub.broadcast()
	server := httptest.NewServer(router)
	defer server.Close()

	// The test server does not use TLS, so we have to convert http to ws.
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Create two clients, one is sender and one is receiver.
	sender, _, err := websocket.DefaultDialer.Dial(wsURL+"/ws?name=sender", nil)
	if err != nil {
		t.Fatalf("dial sender error: %v", err)
	}
	defer sender.Close()

	receiver, _, err := websocket.DefaultDialer.Dial(wsURL+"/ws?name=receiver", nil)
	if err != nil {
		t.Fatalf("dial receiver error: %v", err)
	}
	defer receiver.Close()

	// Sender sends a message.
	msg := InMessage{Content: "hello"}
	bytes, _ := json.Marshal(msg)
	if err := sender.WriteMessage(websocket.TextMessage, bytes); err != nil {
		t.Fatalf("sender write message error: %v", err)
	}

	// Receiver expects to receive a message.
	// It may take a while for the message to be broadcasted, so we have to wait.
	receiver.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, p, err := receiver.ReadMessage()
	if err != nil {
		t.Fatalf("receiver read message error: %v", err)
	}
	var out OutMessage
	if err := json.Unmarshal(p, &out); err != nil {
		t.Fatalf("unmarshal message error: %v", err)
	}
	if out.Content != "hello" {
		t.Fatalf("invalid message content: %v", out.Content)
	}
	if out.Username != "sender" {
		t.Fatalf("invalid message username: %v", out.Username)
	}

	// Sender should not receive any message.
	sender.SetReadDeadline(time.Now().Add(1 * time.Second))
	if _, _, err := sender.ReadMessage(); err == nil {
		t.Fatal("sender should not receive any message")
	}
}

func TestHandleWSMissingName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	hub := NewChatHub()
	router.GET("/ws", hub.handleWS)

	// Create a new request.
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ws", nil)
	router.ServeHTTP(w, req)

	// Check response.
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid status code: %v", w.Code)
	}
	if w.Body.String() != "missing query parameter `name`" {
		t.Fatalf("invalid body: %v", w.Body.String())
	}
}
