package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"message/usecases"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	TYPE = "type"
	AUTH = "auth"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MessageHandler struct {
	au usecases.AuthUsecase
	cu usecases.ChatUsecase
	mb interface{}
}

// 웹소켓은 하나의 핸들러에서 처리 단, useCase를 여러개 둘 수 있음.
func NewMessageHandler(cu usecases.ChatUsecase, au usecases.AuthUsecase, mb interface{}) *MessageHandler {
	return &MessageHandler{
		cu: cu,
		au: au,
		mb: mb,
	}
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
		// 메시지는 반복해서 수신
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read msg error:", err)
		}

		/* 페이로드 검증 */
		var payload map[string]interface{}
		json.Unmarshal(msg, &payload)
		fmt.Println("웹소켓 메시지 수신시 최초 로깅 : ", payload)
		msgType := payload[TYPE].(string)

		// 최초 인증 처리
		if !authenticated {
			if msgType == AUTH {
				token := payload["token"].(string)

				flag, err := h.au.AuthenticateToken(token)
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"auth_failed"}`))
					// defer conn.Close() 호출 되므로 연결 종료.
					break
				}
				authenticated = flag
				conn.WriteMessage(websocket.TextMessage, []byte(`{"success":"auth_ok"}`))
			} else {
				// 인증 안 된 상태에서는 다른 메시지 못 보내게 함
				// defer conn.Close() 호출 되므로 연결 종료.
				conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"unauthorized"}`))
				break
			}
		}

		fmt.Println("인증 완료된 계정 :", payload["user_id"])

		// 여기서 부터 케이스에 따른 처리
		// 채팅방 구독, 채팅방 구독 취소
		// 쪽지

	}

}
