package di

import (
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
	"message/internal/infrastructure/repository"
	"message/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type chatRoomModule struct {
	Handler *handler.ChatRoomHandler
	Useacse usecase.ChatRoomUsecase
}

func InitChatRoomModule(db *gorm.DB, storage storage.ChatRoomStorage) *chatRoomModule {

	repository := repository.NewChatRoomRepository(db)
	usecase := usecase.NewChatRoomUsecase(repository, storage)
	handler := handler.NewChatRoomHandler(usecase)

	return &chatRoomModule{
		Handler: handler,
		Useacse: usecase,
	}
}
