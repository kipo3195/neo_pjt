package usecase

import (
	"message/internal/domain/chat/repository"

	"github.com/nats-io/nats.go"
)

type chatUsecase struct {
	repository repository.ChatRepository
	connector  *nats.Conn
}

type ChatUsecase interface {
	//	SendChat(ctx context.Context, in input.SendChatInput) error
}

func NewChatUsecase(repository repository.ChatRepository, connector *nats.Conn) ChatUsecase {
	return &chatUsecase{
		repository: repository,
		connector:  connector,
	}
}

// func (u *chatUsecase) SendChat(ctx context.Context, in input.SendChatInput) error {

// 	entity := entity.MakeSendChatEntity("", in.SendUserHash, in.Contents, in.LineKey, in.RecvUserHash)
// 	data, err := json.Marshal(entity) // 🔹 struct → []byte(JSON)
// 	if err != nil {
// 		log.Fatal(err)
// 		return err
// 	}

// 	// 채팅 발송
// 	err = u.connector.Publish("chat.message", data)
// 	if err != nil {
// 		log.Fatal("NATS publish failed:", err)
// 		return consts.ErrPublishToMessageBrokerError
// 	}

// 	return nil
// }
