package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"test/dto"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

// roomKey 기반 채널 저장소
var rooms = make(map[string]chan []byte)
var roomsLock sync.Mutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func websocketHandler(w http.ResponseWriter, r *http.Request, nc *nats.Conn) {
	// WebSocket 업그레이드
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// 최초 메시지에서 roomKey를 수신함.
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Initial message read error:", err)
		return
	}

	// 메시지 -> entity
	var req dto.WSRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		log.Println("Invalid JSON:", err)
		return
	}

	roomKey := req.RoomKey
	if roomKey == "" {
		log.Println("Missing roomKey")
		return
	}

	log.Printf("Client joined room: %s\n", roomKey)

	// roomKey에 해당하는 채널 생성 (또는 재사용)
	roomsLock.Lock()
	if _, exists := rooms[roomKey]; !exists {
		rooms[roomKey] = make(chan []byte, 10) // 버퍼는 상황에 맞게 조절
	}
	roomChan := rooms[roomKey]
	roomsLock.Unlock()

	// 메시지 수신 고루틴 → 채널로 전달
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("Received message in room [%s]: %s\n", roomKey, msg)
			roomChan <- msg // 해당 방에 메시지 전달
		}
	}()

	// 채널 메시지를 WebSocket으로 전송
	for msg := range roomChan {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func main() {
	// NATS 서버 연결
	opts := nats.Options{
		Url: "nats://127.0.0.1:4222",
	}

	nc, err := opts.Connect()
	if err != nil {
		log.Fatal("Error connecting to NATS:", err)
	}
	defer nc.Close()

	// 연결된 NATS 서버 URL 출력
	fmt.Println("Connected to NATS at:", opts.Url)

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		websocketHandler(w, r, nc)
	})

	// 서버 실행
	log.Println("WebSocket server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
