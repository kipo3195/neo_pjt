package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/infrastructure/broker"
	"notificator/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type ChatModule struct {
	Usecase usecase.ChatUsecase
}

func InitChatModule(db *gorm.DB, mb broker.Broker) *ChatModule {

	repo := repository.NewChatRepository(db)
	usecase := usecase.NewChatUsecase(repo, mb)

	return &ChatModule{
		Usecase: usecase,
	}
}
