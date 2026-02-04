package di

import (
	"message/internal/adapter/http/handler"
	"message/internal/application/usecase"
	"message/internal/infrastructure/persistence/repository"

	"gorm.io/gorm"
)

type ChatRoomTitleModule struct {
	Handler *handler.ChatRoomTitleHandler
	Usecase usecase.ChatRoomTitleUsecase
}

func InitChatRoomTitleModule(db *gorm.DB) ChatRoomTitleModule {

	repo := repository.NewChatRoomTitleRepository(db)
	usecase := usecase.NewChatRoomTitleUsecase(repo)
	handler := handler.NewChatRoomTitleHandler(usecase)

	return ChatRoomTitleModule{
		Handler: handler,
		Usecase: usecase,
	}
}
