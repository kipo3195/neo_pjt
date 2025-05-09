package usecases

import (
	"encoding/json"
	"fmt"
	"log"
	"message/broker"
	"message/repositories"

	"github.com/gorilla/websocket"
)

const (
	CMD = "cmd" // 이벤트 type의 세부 타입
)

type chatUsecase struct {
	repo repositories.ChatRepository
	mb   broker.Broker
}

type ChatUsecase interface {
	HandleChat(conn *websocket.Conn, data map[string]interface{})
	handleJoinRoom(payload map[string]interface{})
	handleJoinRoomCancle(payload map[string]interface{})
	handleSendMessage(userId string, payload map[string]interface{})
}

func NewChatUsecase(repo repositories.ChatRepository, mb broker.Broker) ChatUsecase {
	return &chatUsecase{
		repo: repo,
		mb:   mb,
	}
}

func (r *chatUsecase) HandleChat(conn *websocket.Conn, data map[string]interface{}) {

	payload := data["payload"].(map[string]interface{})
	cmd := data[CMD].(string)
	userId := data["userId"].(string)

	switch cmd {
	case "joinRoom":
		r.handleJoinRoom(payload)
	case "joinRoomCancel":
		r.handleJoinRoomCancle(payload)
	case "sendMessage":
		r.handleSendMessage(userId, payload)
	}
}

func (r *chatUsecase) handleJoinRoom(payload map[string]interface{}) {
	roomId := payload["roomId"].(string)
	ch, err := r.mb.Subscribe(roomId)

	if err != nil {
		log.Println("subscribe error:", err)
		return
	}

	go func() {
		for msg := range ch {
			fmt.Printf("Received message in room %s: %s\n", roomId, msg)
		}
	}()
}

func (r *chatUsecase) handleJoinRoomCancle(payload map[string]interface{}) {
	roomId := payload["roomId"].(string)
	r.mb.Unsubscribe(roomId)

}

func (r *chatUsecase) handleSendMessage(userId string, payload map[string]interface{}) {
	roomId := payload["roomId"].(string)
	content := payload["content"].(string)
	jsonBytes, _ := json.Marshal(content)

	if err := r.mb.Publish(roomId, jsonBytes); err != nil { // Publish도 하나의 인터페이스에 속한 메소드 구현한다면 Broker의 인터페이스. (덕타이핑)
		log.Println("Publish error:", err)
	}
}
