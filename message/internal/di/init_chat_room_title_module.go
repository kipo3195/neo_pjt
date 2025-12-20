package di

import (
	"message/internal/application/usecase"
	"message/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type ChatRoomTitleModule struct {
	Usecase usecase.ChatRoomTitleUsecase
}

func InitChatRoomTitleModule(db *gorm.DB) ChatRoomTitleModule {

	repo := repository.NewChatRoomTitleRepository(db)
	usecase := usecase.NewChatRoomTitleUsecase(repo)

	return ChatRoomTitleModule{
		Usecase: usecase,
	}
}
