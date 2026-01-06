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

func InitChatModule(db *gorm.DB, chatRoomStorage storage.ChatRoomStorage, sendConnectionStorage storage.SendConnectionStorage) *ChatModule {

	repo := repository.NewChatRepository(db)
	usecase := usecase.NewChatUsecase(chatRoomStorage, sendConnectionStorage, repo)

	return &ChatModule{
		Usecase: usecase,
	}
}
