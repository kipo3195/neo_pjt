package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type ChatModule struct {
	Usecase usecase.ChatUsecase
}

func InitChatModule(db *gorm.DB, chatUserStorage storage.ChatUserStorage) *ChatModule {

	//repo := repository.NewChatRepository(db)
	usecase := usecase.NewChatUsecase(chatUserStorage)

	return &ChatModule{
		Usecase: usecase,
	}
}
