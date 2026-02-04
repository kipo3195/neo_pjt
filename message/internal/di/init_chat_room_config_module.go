package di

import (
	"message/internal/application/usecase"
	"message/internal/infrastructure/persistence/repository"

	"gorm.io/gorm"
)

type ChatRoomConfigModule struct {
	Usecase usecase.ChatRoomConfigUsecase
}

func InitChatRoomConfigModule(db *gorm.DB) ChatRoomConfigModule {

	repo := repository.NewChatRoomConfigRepository(db)
	usecase := usecase.NewChatRoomConfigUsecase(repo)

	return ChatRoomConfigModule{
		Usecase: usecase,
	}
}
