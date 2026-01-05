package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/infrastructure/repository"
	"notificator/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type ChatModule struct {
	Usecase usecase.ChatUsecase
}

func InitChatModule(db *gorm.DB, chatUserStorage storage.ChatUserStorage, sendConnectionStorage storage.SendConnectionStorage) *ChatModule {

	repo := repository.NewChatRepository(db)
	usecase := usecase.NewChatUsecase(chatUserStorage, sendConnectionStorage, repo)

	return &ChatModule{
		Usecase: usecase,
	}
}
