package nats

import (
	"context"
	"encoding/json"
	"log"
	"notificator/internal/application/usecase"
	"notificator/internal/application/usecase/input"
	"notificator/internal/delivery/adapter"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsChatSubscriber struct {
	conn                *nats.Conn
	chatUsecase         usecase.ChatUsecase
	socketSenderUsecase usecase.SocketSenderUsecase
	handler             func(data []byte) error
}

func NewNatsChatSubscriber(nc *nats.Conn, chatUsecase usecase.ChatUsecase, socketSendUsecase usecase.SocketSenderUsecase) *NatsChatSubscriber {
	return &NatsChatSubscriber{
		conn:                nc,
		chatUsecase:         chatUsecase,
		socketSenderUsecase: socketSendUsecase,
	}
}

// 구독, goroutine을 동한 처리 변경
func (s *NatsChatSubscriber) AddSubscribe(kind string) error {

	// NATS로부터 메시지를 하나 받을 때마다 go s.handleNatsMessage(kind, msg.Data)를 호출
	_, err := s.conn.Subscribe(kind, func(msg *nats.Msg) {
		// 수신 데이터 로깅
		log.Println("[Notificator] kind : "+kind+" Received message:", string(msg.Data))

		// 수신 받은 데이터는 별도 고루틴에서 처리
		go s.handleNatsMessage(kind, msg.Data)
	})

	if err != nil {
		log.Printf("NATS subscription failed for %s: %v", kind, err)
		return err
	}

	return nil
}

func (s *NatsChatSubscriber) handleNatsMessage(kind string, data []byte) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch kind {
	case "chat.broadcast":

		var input input.ChatMessageInput
		if err := json.Unmarshal(data, &input); err != nil {
			log.Printf("invalid message: %v", err)
			return
		}

		// chat 메시지 -> 클라이언트
		output := s.chatUsecase.RecvChatMessage(ctx, input)
		// 실시간 발송 처리를 위한 도메인 구분 (chatUsecase, socketSenderUsecase)
		in := adapter.MakeChatInput(output.EventType, output.ChatSession, output.ChatRoomData, output.ChatLineData)
		s.socketSenderUsecase.RecvChat(ctx, in)

	}
}
