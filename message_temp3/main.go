package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

type WSRequest struct {
	RoomKey string `json:"roomKey"`
	UserId  string `json:"userId"`
	Type    string `json:"type"`
}

type incomming struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("NATS connection error:", err)
	}
	defer nc.Drain()

	http.HandleFunc("/ws", websocketHandler(nc))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func websocketHandler(nc *nats.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}
		defer conn.Close()

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read roomKey error:", err)
			return
		}

		var req WSRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Println("Invalid roomKey JSON:", err)
			return
		}

		roomKey := req.RoomKey
		if roomKey == "" {
			log.Println("Missing roomKey")
			return
		}

		log.Printf("Client joined room: %s", roomKey)

		// 채널 생성
		msgChan := make(chan *nats.Msg, 64) // 버퍼 크기 설정

		sub, err := nc.ChanSubscribe(roomKey, msgChan)
		if err != nil {
			log.Println("NATS subscribe error:", err)
			return
		}
		defer sub.Unsubscribe()

		// 메인 select 루프
		for {
			select {
			// NATS → WebSocket
			case m := <-msgChan:
				var initReq incomming
				_ = json.Unmarshal(m.Data, &initReq)
				senderID := initReq.Sender
				if req.UserId == senderID {
					continue // 내 메시지는 무시
				}
				if err := conn.WriteMessage(websocket.TextMessage, m.Data); err != nil {
					log.Println("WebSocket write error:", err)
					return // 에러 발생 시 루프 종료 (연결 끊김)
				}

			// WebSocket → NATS
			default:
				// Non-blocking read
				conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond)) // 짧은 타임아웃
				_, wsMsg, err := conn.ReadMessage()
				fmt.Println("여기 출력 :", string(wsMsg))
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Println("WebSocket read error:", err)
					}
					return // 정상 종료 (클라이언트가 끊음)
				}

				if err := nc.Publish(roomKey, wsMsg); err != nil {
					log.Println("NATS publish error:", err)
				}
			}
		}
	}
}
