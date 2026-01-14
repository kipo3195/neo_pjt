package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/domain/port"
	"notificator/internal/infrastructure/repository"

	"notificator/internal/infrastructure/storage"
	"notificator/internal/infrastructure/workerPool"

	"gorm.io/gorm"
)

type ChatModule struct {
	Usecase usecase.ChatUsecase
}

func InitChatModule(db *gorm.DB, chatRoomStorage storage.ChatRoomStorage, messageSender port.MessageSender) *ChatModule {

	repo := repository.NewChatRepository(db)

	workerPool := workerPool.NewChatWorkerPool(10, messageSender)
	workerPool.Init()
	usecase := usecase.NewChatUsecase(chatRoomStorage, repo, messageSender, workerPool)

	return &ChatModule{
		Usecase: usecase,
	}
}
