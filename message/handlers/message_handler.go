package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"message/broker"
	"message/usecases"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	TYPE = "type" // 이벤트의 대분류
	AUTH = "auth"
	CHAT = "chat"
	NOTE = "note"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MessageHandler struct {
	au usecases.AuthUsecase
	cu usecases.ChatUsecase
	nu usecases.NoteUsecase
}

// 웹소켓은 하나의 핸들러에서 처리 단, useCase를 여러개 둘 수 있음.
func NewMessageHandler(cu usecases.ChatUsecase, au usecases.AuthUsecase, nu usecases.NoteUsecase, mb broker.Broker) *MessageHandler {
	return &MessageHandler{
		cu: cu,
		au: au,
		nu: nu,
	}
}

type incomming struct {
	Content string `json:"content"`
	Sender  string `json:"sender"`
}

// MessageHandler의 메소드 SetupRoutes에서 등록됨.
func (h *MessageHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// 인증 처리 flag
	authenticated := false

	for {
		// 메시지는 반복해서 수신, ReadMessage는 블로킹 함수
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read msg error:", err)
		}

		/* 페이로드 검증 */
		var data map[string]interface{}
		json.Unmarshal(msg, &data)
		fmt.Println("웹소켓 메시지 수신시 최초 로깅 : ", data)
		msgType, ok := data[TYPE].(string)
		if !ok {
			conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"invalid_type"}`))
			continue
		}

		/*최초 인증 처리*/
		if !authenticated {

			switch msgType {
			case AUTH:
				result, err := h.au.HandleAuth(conn, data)
				if !result && err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"auth_failed"}`))
					return
				} else {
					conn.WriteMessage(websocket.TextMessage, []byte(`{"success":"auth_ok"}`))
					authenticated = true
					continue
				}
			default:
				conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"unauthorized"}`))
				return
			}
		}

		// 인증 이후에는 웹소켓 수신 메시지를 가지고 이쪽으로
		if authenticated {
			switch msgType {
			case CHAT:
				h.cu.HandleChat(conn, data)
			case NOTE:
				h.nu.HandleNote(conn, data)
			}
		}

	}

}
