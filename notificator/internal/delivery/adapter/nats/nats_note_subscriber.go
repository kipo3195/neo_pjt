package nats

import (
	"encoding/json"
	"log"
	"notificator/internal/application/usecase"
	"notificator/internal/application/usecase/input"

	"github.com/nats-io/nats.go"
)

type NatsNoteSubscriber struct {
	conn                *nats.Conn
	noteUsecase         usecase.NoteUsecase
	socketSenderUsecase usecase.SocketSenderUsecase
	handler             func(data []byte) error
}

func NewNatsNoteSubscriber(nc *nats.Conn, noteUsecase usecase.NoteUsecase, socketSendUsecase usecase.SocketSenderUsecase) *NatsNoteSubscriber {
	return &NatsNoteSubscriber{
		conn:                nc,
		noteUsecase:         noteUsecase,
		socketSenderUsecase: socketSendUsecase,
	}
}

// 구독, goroutine을 동한 처리 변경
func (s *NatsNoteSubscriber) AddSubscribe(kind string) error {

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

func (s *NatsNoteSubscriber) handleNatsMessage(kind string, data []byte) {

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	switch kind {

	case "note.broadcast":
		var input input.NoteMessageInput
		if err := json.Unmarshal(data, &input); err != nil {
			log.Printf("invalid message: %v", err)
			return
		}
	}
}
