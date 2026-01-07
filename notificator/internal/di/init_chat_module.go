package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/core/port"
	"notificator/internal/infrastructure/repository"
	"notificator/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type ChatModule struct {
	Usecase usecase.ChatUsecase
}

func InitChatModule(db *gorm.DB, chatRoomStorage storage.ChatRoomStorage, messageSender port.MessageSender) *ChatModule {

	repo := repository.NewChatRepository(db)
	usecase := usecase.NewChatUsecase(chatRoomStorage, repo, messageSender)

	return &ChatModule{
		Usecase: usecase,
	}
}
