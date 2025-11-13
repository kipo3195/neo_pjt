package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"notificator/internal/application/orchestrator"
	"notificator/internal/consts"
	"notificator/internal/delivery/dto/chat"
	"notificator/internal/delivery/dto/notificatorService"

	"github.com/gorilla/websocket"
)

type NotificatorServiceHandler struct {
	svc *orchestrator.NotificatorService
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewNotificatorServiceHandler(svc *orchestrator.NotificatorService) *NotificatorServiceHandler {
	return &NotificatorServiceHandler{
		svc: svc,
	}
}

func (h *NotificatorServiceHandler) NotificatorConnect(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Notificator service WebSocket upgrade error:", err)
		return
	}

	defer conn.Close()

	for {
		// 메시지는 반복해서 수신, ReadMessage는 블로킹 함수
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Notificator service Read msg error:", err)
			break
		}

		// type 파싱
		var req notificatorService.NotificatorConnectRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Println("Notificator service websocket message error:", err)
			return
		}

		// 여기서 각각의 usecase를 활용한 처리

		switch req.Type {

		case consts.AUTH:

		case consts.CHAT:
			var chatMessage chat.ChatConnect
			if err := json.Unmarshal(msg, &chatMessage); err == nil {
				h.svc.Chat.SubscribeChat(chatMessage, conn)
			}

		case consts.NOTE:

		default:
			log.Println("unknown message type:", req.Type)
			return
		}

	}

	log.Println("Notificator service websocket close ")

}
