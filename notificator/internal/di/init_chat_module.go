package di

import (
	"log"
	"notificator/internal/application/usecase"
	"notificator/internal/domain/port"
	"notificator/internal/infrastructure/repository"

	"notificator/internal/infrastructure/storage"
	"notificator/internal/infrastructure/workerPool"

	"gorm.io/gorm"
)

type ChatModule struct {
	Usecase    usecase.ChatUsecase
	workerPool workerPool.ChatWorkerPool
}

func InitChatModule(db *gorm.DB, chatRoomStorage storage.ChatRoomStorage, messageSender port.MessageSender) *ChatModule {

	repo := repository.NewChatRepository(db)

	workerPool := workerPool.NewChatWorkerPool(10, messageSender)
	workerPool.Init()
	usecase := usecase.NewChatUsecase(chatRoomStorage, repo, messageSender, workerPool)

	return &ChatModule{
		Usecase:    usecase,
		workerPool: workerPool,
	}
}

func (m *ChatModule) Cleanup() {
	log.Println("Cleaning up ChatModule...")
	// WorkerPool 종료 (채널 닫기 및 대기)
	m.workerPool.Stop()

	// 추가적으로 NATS 연결 종료 등이 필요하다면 여기서 수행
	// m.Connector.Close()
}
