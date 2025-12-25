package di

import (
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
	"message/internal/infrastructure/repository"

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
