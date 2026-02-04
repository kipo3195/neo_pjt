package di

import (
	"message/internal/application/usecase"
	"message/internal/infrastructure/persistence/repository"

	"gorm.io/gorm"
)

type ChatFileModule struct {
	Usecase usecase.ChatFileUsecase
}

func InitChatFileModule(db *gorm.DB) *ChatFileModule {

	repository := repository.NewChatFileRepository(db)
	usecase := usecase.NewChatFileUsecase(repository)

	return &ChatFileModule{
		Usecase: usecase,
	}

}
