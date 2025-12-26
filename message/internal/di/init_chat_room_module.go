package di

import (
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
	"message/internal/infrastructure/repository"
	"message/internal/infrastructure/storage"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type chatRoomModule struct {
	Handler *handler.ChatRoomHandler
	Usecase usecase.ChatRoomUsecase
}

func InitChatRoomModule(db *gorm.DB, storage storage.ChatRoomStorage, connector *nats.Conn) *chatRoomModule {

	repository := repository.NewChatRoomRepository(db)
	usecase := usecase.NewChatRoomUsecase(repository, connector, storage)
	handler := handler.NewChatRoomHandler(usecase)

	return &chatRoomModule{
		Handler: handler,
		Usecase: usecase,
	}
}
