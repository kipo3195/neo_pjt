package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

type WSRequest struct {
	RoomKey string `json:"roomKey"`
	UserId  string `json:"userId"`
}

type incomming struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
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
		// WebSocket 업그레이드
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}
		defer conn.Close()

		// 최초 메시지에서 roomKey 추출
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read roomKey error:", err)
			return
		}

		// json 형태의 데이터 전달
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

		// Subscribe to NATS roomKey
		// Subscribe는 NATS 클라이언트 라이브러리의 내부 구현이 콜백 실행을 고루틴으로 처리한다.
		// 즉, 비동기로 대기중인 상태이므로 종료되지 않는다.
		//  func(m *nats.Msg) 가 콜백함수

		sub, err := nc.Subscribe(roomKey, func(m *nats.Msg) {
			// sender가 보낸 처리를 막음
			var initReq incomming
			_ = json.Unmarshal(m.Data, &initReq)
			senderID := initReq.Sender
			fmt.Printf("senderID : %s, req.UserId : %s \n", senderID, req.UserId)
			if req.UserId == senderID {
				return
			}
			conn.WriteMessage(websocket.TextMessage, m.Data)
		})

		if err != nil {
			log.Println("NATS subscribe error:", err)
			return
		}
		defer sub.Unsubscribe()

		// WebSocket → NATS publisher
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("WebSocket read error:", err)
				break
			}

			log.Printf("Publishing to room %s: %s", roomKey, string(msg))
			// json 형태로 발송
			if err := nc.Publish(roomKey, msg); err != nil {
				log.Println("NATS publish error:", err)
			}
		}
	}
}
