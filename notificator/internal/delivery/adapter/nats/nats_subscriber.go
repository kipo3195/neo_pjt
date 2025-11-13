package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notificator/internal/application/usecase"
	"notificator/internal/delivery/adapter"
	"notificator/internal/delivery/dto/chat"

	"github.com/nats-io/nats.go"
)

type NatsSubscriber struct {
	conn        *nats.Conn
	chatUsecase usecase.ChatUsecase
	handler     func(data []byte) error
}

func NewNatsSubscriber(nc *nats.Conn, chatUsecase usecase.ChatUsecase) *NatsSubscriber {
	return &NatsSubscriber{conn: nc, chatUsecase: chatUsecase}
}

// 구독 후 메시지 수신
func (s *NatsSubscriber) Subscribe(ctx context.Context) error {
	_, err := s.conn.Subscribe("chat.message", func(msg *nats.Msg) {
		var data chat.ChatMessage
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("invalid message: %v", err)
			return
		}

		input := adapter.MakeChatMessageInput(data.Type, data.SendUserHash, data.Contents, data.DestUserHash)
		s.chatUsecase.RecvChatMessage(ctx, input)
		// if err := s.chatUsecase.RecvChatMessage(ctx, input); err != nil {
		// 	log.Printf("usecase error: %v", err)
		// }
	})
	return err
}

// 구독
func (s *NatsSubscriber) StartSubscribe(kind string) error {
	_, err := s.conn.Subscribe(kind, func(msg *nats.Msg) {
		if err := s.handler(msg.Data); err != nil {
			fmt.Println("NATS handle error:", err)
		}
	})
	return err
}
