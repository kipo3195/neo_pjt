package di

import (
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
	"message/internal/infrastructure/repository"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type ChatModule struct {
	Handler *handler.ChatHandler
	Usecase usecase.ChatUsecase
}

func InitChatModule(db *gorm.DB, connector *nats.Conn) *ChatModule {
	repository := repository.NewChatRepository(db)
	usecase := usecase.NewChatUsecase(repository, connector)
	handler := handler.NewChatHandler(usecase)

	return &ChatModule{
		Handler: handler,
		Usecase: usecase,
	}
}
