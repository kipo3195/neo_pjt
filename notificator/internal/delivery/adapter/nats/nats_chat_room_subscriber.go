package nats

import (
	"context"
	"encoding/json"
	"log"
	"notificator/internal/application/usecase"
	"notificator/internal/application/usecase/input"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsChatRoomSubscriber struct {
	conn                *nats.Conn
	chatUsecase         usecase.ChatUsecase
	noteUsecase         usecase.NoteUsecase
	socketSenderUsecase usecase.SocketSenderUsecase
	handler             func(data []byte) error
}

func NewNatsChatRoomSubscriber(nc *nats.Conn, chatUsecase usecase.ChatUsecase, noteUsecase usecase.NoteUsecase, socketSendUsecase usecase.SocketSenderUsecase) *NatsChatRoomSubscriber {
	return &NatsChatRoomSubscriber{conn: nc, chatUsecase: chatUsecase, noteUsecase: noteUsecase, socketSenderUsecase: socketSendUsecase}
}

// 구독, goroutine을 동한 처리 변경
func (s *NatsChatRoomSubscriber) AddSubscribe(kind string) error {

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

func (s *NatsChatRoomSubscriber) handleNatsMessage(kind string, data []byte) {

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	switch kind {
	case "chat.room.broadcast":

	}
}

// 로드밸런싱 정책이 Queue Group (구독하는 N개중 1개만 수신 할 수 있음. )
func (s *NatsChatRoomSubscriber) AddQueueSubscribe(kind string) error {

	_, err := s.conn.QueueSubscribe(kind, "notificator-queue-group", func(msg *nats.Msg) {
		// 수신 데이터 로깅
		log.Println("[Notificator Queue Group] kind : "+kind+" Received message:", string(msg.Data))

		// 수신 받은 데이터는 별도 고루틴에서 처리
		go s.QueueGrouphandleNatsMessage(kind, msg.Data)

		// 처리가 끝났음을 알림 (Reply)
		msg.Respond([]byte(kind + " success"))
	})

	if err != nil {
		log.Printf("NATS subscription failed for %s: %v", kind, err)
		return err
	}

	return nil
}

func (s *NatsChatRoomSubscriber) QueueGrouphandleNatsMessage(kind string, data []byte) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch kind {
	case "create.chat.room":
		var input input.CreateChatRoomMessageInput
		if err := json.Unmarshal(data, &input); err != nil {
			log.Printf("invalid message: %v", err)
			return
		}
		log.Printf("[%s] input : %s\n", kind, input)
		// 실시간 발송 처리를 위한 도메인 구분 (chatUsecase, socketSenderUsecase)
		err := s.chatUsecase.RecvCreateChatRoomMessage(ctx, input)
		if err != nil {
			log.Printf("[%s] err : %s\n", kind, err)
			return
		}

		// 별도의 가공처리가 필요 없음, RecvCreateChatRoomMessage에서도 별도의 가공처리를 하지 않으므로 input을 그대로 사용함.
		s.socketSenderUsecase.RecvCreateChatRoom(ctx, input)

		log.Printf("[%s] success. \n", kind)
	}
}
