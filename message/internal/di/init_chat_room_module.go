package di

import (
	"message/internal/adapter/http/handler"
	"message/internal/application/usecase"
	"message/internal/domain/logger"
	"message/internal/infrastructure/persistence/repository"
	"message/internal/infrastructure/storage"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type chatRoomModule struct {
	Handler *handler.ChatRoomHandler
	Usecase usecase.ChatRoomUsecase
}

func InitChatRoomModule(db *gorm.DB, storage storage.ChatRoomStorage, connector *nats.Conn, logger logger.Logger) *chatRoomModule {

	repository := repository.NewChatRoomRepository(db)
	usecase := usecase.NewChatRoomUsecase(repository, connector, storage, logger)
	handler := handler.NewChatRoomHandler(usecase)

	return &chatRoomModule{
		Handler: handler,
		Usecase: usecase,
	}
}
