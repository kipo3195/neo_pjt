package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/domain/port"
	"notificator/internal/infrastructure/repository"
	"notificator/internal/infrastructure/storage"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type ChatRoomModule struct {
	Usecase usecase.ChatRoomUsecase
}

func InitChatRoomModule(db *gorm.DB, chatRoomStorage storage.ChatRoomStorage, sendConnectionStorage storage.SendConnectionStorage, connector *nats.Conn, messageSender port.MessageSender) *ChatRoomModule {

	repo := repository.NewChatRoomRepository(db)
	usecase := usecase.NewChatRoomUsecase(repo, chatRoomStorage, sendConnectionStorage, connector, messageSender)

	return &ChatRoomModule{
		Usecase: usecase,
	}
}
