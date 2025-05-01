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

const (
	// ping/pong 타임아웃 기준
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10 // ping은 조금 더 자주 보냄
)

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

		// conn.SetPongHandler() → 클라이언트가 pong 응답을 보내면 read deadline을 자동 연장.
		conn.SetPongHandler(func(appData string) error {
			log.Println("Pong received")
			return conn.SetReadDeadline(time.Now().Add(pongWait))
		})

		// ✅ 최초 ReadDeadline 설정
		conn.SetReadDeadline(time.Now().Add(pongWait))

		// ✅ ping ticker 시작
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		// 최초 메시지에서 roomKey 추출
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
		msgChan := make(chan *nats.Msg, 64)

		// NATS 서버에 roomKey라는 subject(=토픽)에 대한 구독을 설정하고, 이 subject로 발행되는 메시지를 msgChan이라는 Go 채널로 전달받게 합니다
		sub, err := nc.ChanSubscribe(roomKey, msgChan)
		if err != nil {
			log.Println("NATS subscribe error:", err)
			return
		}
		defer sub.Unsubscribe()

		for {
			select {
			case m := <-msgChan:
				var initReq incomming
				_ = json.Unmarshal(m.Data, &initReq)
				fmt.Printf("채널에서 채팅 수신 ! sender : %s , message : %s \n", initReq.Sender, initReq.Message)
				senderID := initReq.Sender
				if req.UserId == senderID {
					continue
				}
				conn.WriteMessage(websocket.TextMessage, m.Data)

			case <-ticker.C:
				// ticker로 주기적으로 PingMessage를 보냄 → 클라이언트가 살아있으면 pong 응답.
				log.Println("Sending ping")
				conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Println("Ping error:", err)
					return
				}

			default:
				// SetReadDeadline()은 pongWait 주기로 계속 갱신됨 → 타임아웃은 죽은 연결에만 발생.
				conn.SetReadDeadline(time.Now().Add(pongWait))
				// 웹소켓에서 메시지를 읽는 부분
				_, msg, err := conn.ReadMessage() // conn.ReadMessage()는 블록킹 함수임. 블록킹은
				fmt.Println("웹소켓 데이터 수신 ! msg : ", string(msg))
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("Unexpected close error: %v", err)
					} else {
						log.Println("Read error:", err)
					}
					return
				}
				fmt.Println("웹소켓 데이터 수신한 데이터를 room에 발송 : ", roomKey)
				if err := nc.Publish(roomKey, msg); err != nil {
					log.Println("NATS publish error:", err)
				}
			}
		}
	}
}
