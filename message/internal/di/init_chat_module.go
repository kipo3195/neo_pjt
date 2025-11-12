package di

import (
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
	"message/internal/infrastructure/broker"
	"message/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type ChatModule struct {
	Handler *handler.ChatHandler
	Usecase usecase.ChatUsecase
}

func InitChatModule(db *gorm.DB, mb broker.Broker) *ChatModule {
	repository := repository.NewChatRepository(db)
	usecase := usecase.NewChatUsecase(repository, mb)
	handler := handler.NewChatHandler(usecase)

	return &ChatModule{
		Handler: handler,
		Usecase: usecase,
	}
}
