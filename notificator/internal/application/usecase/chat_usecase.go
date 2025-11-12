package usecase

import (
	"notificator/internal/delivery/dto/chat"
	"notificator/internal/domain/chat/repository"
	"notificator/internal/infrastructure/broker"

	"github.com/gorilla/websocket"
)

const (
	CMD = "cmd" // 이벤트 type의 세부 타입
)

type chatUsecase struct {
	repo repository.ChatRepository
	mb   broker.Broker
}

type ChatUsecase interface {
	// HandleChat(conn *websocket.Conn, data map[string]interface{})
	// handleJoinRoom(payload map[string]interface{}, conn *websocket.Conn)
	// handleJoinRoomCancle(payload map[string]interface{})
	// handleSendMessage(userId string, payload map[string]interface{})
	SubscribeChat(chatMessage chat.ChatMessage, conn *websocket.Conn)
}

func NewChatUsecase(repo repository.ChatRepository, mb broker.Broker) ChatUsecase {
	return &chatUsecase{
		repo: repo,
		mb:   mb,
	}
}

func (u *chatUsecase) SubscribeChat(chatMessage chat.ChatMessage, conn *websocket.Conn) {

	// chat + 사용자 hash으로 구독
	u.mb.SubscribeChat(chatMessage.UserHash, conn)

}

// func (r *chatUsecase) HandleChat(conn *websocket.Conn, data map[string]interface{}) {

// 	payload := data["payload"].(map[string]interface{})
// 	cmd := data[CMD].(string)
// 	userId := data["userId"].(string)

// 	switch cmd {
// 	case consts.JoinRoom:
// 		r.handleJoinRoom(payload, conn)
// 	case consts.JoinRoomCancle:
// 		r.handleJoinRoomCancle(payload)
// 	case consts.SendMessage:
// 		r.handleSendMessage(userId, payload)
// 	}
// }

// func (r *chatUsecase) handleJoinRoom(payload map[string]interface{}, conn *websocket.Conn) {
// 	// SoC, usecase에서는 '구독' 만 처리, broker에서 구독에 대한 상세 처리.
// 	roomId := payload[consts.RoomId].(string)
// 	userId := payload[consts.UserId].(string)
// 	r.mb.SubscribeChatRoom(roomId, userId, conn)
// }

// func (r *chatUsecase) handleJoinRoomCancle(payload map[string]interface{}) {
// 	roomId := payload[consts.RoomId].(string)
// 	userId := payload[consts.UserId].(string)
// 	r.mb.UnSubscribeChatRoom(roomId, userId)

// }

// func (r *chatUsecase) handleSendMessage(userId string, payload map[string]interface{}) {
// 	roomId := payload[consts.RoomId].(string)
// 	content := payload[consts.Content].(string)
// 	jsonBytes, _ := json.Marshal(content)

// 	if err := r.mb.PublishToChatRoom(roomId, jsonBytes); err != nil { // Publish도 하나의 인터페이스에 속한 메소드 구현한다면 Broker의 인터페이스. (덕타이핑)
// 		log.Println("Publish error:", err)
// 	}
// }
