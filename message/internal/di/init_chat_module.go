package di

import (
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
	"message/internal/infrastructure/repository"
	"message/internal/infrastructure/workerPool"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type ChatModule struct {
	Handler    *handler.ChatHandler
	Usecase    usecase.ChatUsecase
	WorkerPool workerPool.ChatWorkerPool
}

func InitChatModule(db *gorm.DB, connector *nats.Conn) *ChatModule {

	repository := repository.NewChatRepository(db)

	// 이 영역에서 구현체를 생성하고 인터페이스 타입으로 Usecase에 주입합니다.
	workerPool := workerPool.NewChatWorkerPool(10, repository)
	workerPool.Init()
	usecase := usecase.NewChatUsecase(repository, connector, workerPool)
	handler := handler.NewChatHandler(usecase)

	return &ChatModule{
		Handler:    handler,
		Usecase:    usecase,
		WorkerPool: workerPool,
	}
}
