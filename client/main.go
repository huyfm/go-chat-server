package main

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const NumRepeats = 4

// Represent incoming messages from users to the hub.
type InMessage struct {
	Content string `json:"content"`
}

// Represent outgoing messages broadcasting from the hub.
type OutMessage struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

// defer func() {
// 	// Send close frame first
// 	_ = conn.WriteMessage(websocket.CloseMessage,
// 		websocket.FormatCloseMessage(closeCode, closeMsg))
// 	time.Sleep(1 * time.Second) // optional wait for peer response
// 	_ = conn.Close()
// }()

func sendHello() {
	url := "ws://localhost:5001/ws?name=huy"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("huy: dial error:", err)
	}
	defer conn.Close()

	hello := InMessage{"Hello"}
	for range NumRepeats {
		if err := conn.WriteJSON(&hello); err != nil {
			log.Fatal("huy: send error:", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	log.Println("huy: send done")
}

func recvHello() {
	url := "ws://localhost:5001/ws?name=peter"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("recv: dial error:", err)
	}
	defer conn.Close()

	for range NumRepeats {
		var msg OutMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Fatal("peter: recv error:", err)
		}
		log.Println("peter: recv", msg)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		sendHello()
	}()
	go func() {
		defer wg.Done()
		recvHello()
	}()

	wg.Wait()
}
