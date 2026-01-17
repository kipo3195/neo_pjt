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

type NatsServiceUsersSubscriber struct {
	conn                *nats.Conn
	serviceUsersUsecase usecase.ServiceUsersUsecase
}

func NewNatsServiceUsersSubscriber(nc *nats.Conn, serviceUsersUsecase usecase.ServiceUsersUsecase) *NatsServiceUsersSubscriber {
	return &NatsServiceUsersSubscriber{
		conn:                nc,
		serviceUsersUsecase: serviceUsersUsecase,
	}
}

func (s *NatsServiceUsersSubscriber) AddQueueSubscribe(kind string) error {

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

func (s *NatsServiceUsersSubscriber) QueueGrouphandleNatsMessage(kind string, data []byte) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch kind {
	case "users.registered":
		var input input.ServiceUsersMessageInput
		if err := json.Unmarshal(data, &input); err != nil {
			log.Printf("invalid message: %v", err)
			return
		}
		log.Printf("[%s] input : %s\n", kind, input)
		// 실시간 발송 처리를 위한 도메인 구분 (chatUsecase, socketSenderUsecase)
		err := s.serviceUsersUsecase.RecvRegistServiceUsersMessage(ctx, input)
		if err != nil {
			log.Printf("[%s] err : %s\n", kind, err)
			return
		}

		log.Printf("[%s] success. \n", kind)
	}
}
