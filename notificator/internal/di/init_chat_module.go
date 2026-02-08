package di

import (
	"log"
	"notificator/internal/application/usecase"
	"notificator/internal/domain/port"
	"notificator/internal/infrastructure/persistence/repository"

	"notificator/internal/infrastructure/storage"
	"notificator/internal/infrastructure/workerPool"

	"gorm.io/gorm"
)

type ChatModule struct {
	Usecase             usecase.ChatUsecase
	chatCountWorkerPool workerPool.ChatCountWorkerPool
}

func InitChatModule(db *gorm.DB, chatRoomStorage storage.ChatRoomStorage, messageSender port.MessageSender) *ChatModule {

	repo := repository.NewChatRepository(db)

	chatCountWorkerPool := workerPool.NewChatCountWorkerPool(10, messageSender)
	chatCountWorkerPool.Init()

	chatReadDateWorkerPool := workerPool.NewChatReadDateWorkerPool(10, messageSender)
	chatReadDateWorkerPool.Init()
	usecase := usecase.NewChatUsecase(chatRoomStorage, repo, messageSender, chatCountWorkerPool, chatReadDateWorkerPool)

	return &ChatModule{
		Usecase:             usecase,
		chatCountWorkerPool: chatCountWorkerPool,
	}
}

func (m *ChatModule) Cleanup() {
	log.Println("Cleaning up ChatModule...")
	// WorkerPool 종료 (채널 닫기 및 대기)
	m.chatCountWorkerPool.Stop()

	// 추가적으로 NATS 연결 종료 등이 필요하다면 여기서 수행
	// m.Connector.Close()
}
