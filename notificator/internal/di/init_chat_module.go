package di

import (
	"notificator/internal/application/usecase"

	"gorm.io/gorm"
)

type ChatModule struct {
	Usecase usecase.ChatUsecase
}

func InitChatModule(db *gorm.DB) *ChatModule {

	//repo := repository.NewChatRepository(db)
	usecase := usecase.NewChatUsecase()

	return &ChatModule{
		Usecase: usecase,
	}
}
