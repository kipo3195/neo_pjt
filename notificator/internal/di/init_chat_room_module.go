package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/infrastructure/repository"
	"notificator/internal/infrastructure/storage"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type ChatRoomModule struct {
	Usecase usecase.ChatRoomUsecase
}

func InitChatRoomModule(db *gorm.DB, chatRoomStorage storage.ChatRoomStorage, connector *nats.Conn) *ChatRoomModule {

	repo := repository.NewChatRoomRepository(db)
	usecase := usecase.NewChatRoomUsecase(repo, chatRoomStorage, connector)

	return &ChatRoomModule{
		Usecase: usecase,
	}
}
