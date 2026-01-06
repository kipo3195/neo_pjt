package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/infrastructure/repository"
	"notificator/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type ChatRoomModule struct {
	Usecase usecase.ChatRoomUsecase
}

func InitChatRoomModule(db *gorm.DB, chatRoomStorage storage.ChatRoomStorage) *ChatRoomModule {

	repo := repository.NewChatRoomRepository(db)
	usecase := usecase.NewChatRoomUsecase(repo, chatRoomStorage)

	return &ChatRoomModule{
		Usecase: usecase,
	}
}
