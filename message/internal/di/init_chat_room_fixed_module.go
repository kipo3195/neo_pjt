package di

import (
	"message/internal/application/usecase"
	"message/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type ChatRoomFixedModule struct {
	Usecase usecase.ChatRoomFixedUsecase
}

func InitChatRoomFixedModule(db *gorm.DB) ChatRoomFixedModule {

	repo := repository.NewChatRoomFixedRepository(db)
	usecase := usecase.NewChatRoomFixedUsecase(repo)

	return ChatRoomFixedModule{
		Usecase: usecase,
	}
}
