package nats

import (
	"context"
	"encoding/json"
	"log"
	"notificator/internal/application/usecase"
	"notificator/internal/delivery/adapter"
	"notificator/internal/delivery/dto/chat"
	"time"

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
	// _, err := s.conn.Subscribe("chat.message", func(msg *nats.Msg) {
	// 	var data chat.ChatMessage
	// 	if err := json.Unmarshal(msg.Data, &data); err != nil {
	// 		log.Printf("invalid message: %v", err)
	// 		return
	// 	}

	// 	input := adapter.MakeChatMessageInput(data.Type, data.SendUserHash, data.Contents, data.LineKey, data.DestUserHash)
	// 	s.chatUsecase.RecvChatMessage(ctx, input)
	// 	// if err := s.chatUsecase.RecvChatMessage(ctx, input); err != nil {
	// 	// 	log.Printf("usecase error: %v", err)
	// 	// }
	// })
	// return err
	return nil
}

// 구독, goroutine을 동한 처리 변경
func (s *NatsSubscriber) AddSubscribe(kind string) error {

	sub, err := s.conn.SubscribeSync(kind)
	if err != nil {
		log.Fatal(err)
	}

	// 이후 kind에 따른 분기처리 필요

	go func() {
		for {
			msg, err := sub.NextMsg(time.Hour * 24)
			if err != nil {
				log.Println("NATS receive error:", err)
				continue
			}

			log.Println("Received chat message:", string(msg.Data))

			var data chat.ChatMessage
			if err := json.Unmarshal(msg.Data, &data); err != nil {
				log.Printf("invalid message: %v", err)
				continue
			}

			input := adapter.MakeChatMessageInput(
				data.Type, data.SendUserHash, data.Contents, data.LineKey, data.DestUsers,
			)

			// chat 메시지 -> 클라이언트
			s.chatUsecase.RecvChatMessage(context.Background(), input)
		}
	}()
	return nil
}
