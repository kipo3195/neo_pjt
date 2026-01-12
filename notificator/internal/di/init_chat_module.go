package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/domain/port"
	"notificator/internal/infrastructure/repository"
	"notificator/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type ChatModule struct {
	Usecase usecase.ChatUsecase
}

func InitChatModule(db *gorm.DB, chatRoomStorage storage.ChatRoomStorage, messageSender port.MessageSender, chatDebouncer port.ChatCountDebouncer) *ChatModule {

	repo := repository.NewChatRepository(db)
	usecase := usecase.NewChatUsecase(chatRoomStorage, repo, messageSender, chatDebouncer)

	return &ChatModule{
		Usecase: usecase,
	}
}
