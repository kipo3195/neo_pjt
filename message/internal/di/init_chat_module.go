package di

import (
	"log"
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
	"message/internal/infrastructure/repository"
	"message/internal/infrastructure/workerPool"
	"time"

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

func (m *ChatModule) Cleanup() {
	log.Println("Cleaning up ChatModule...")

	// 채널 닫기 및 대기를 별도 채널로 감시
	done := make(chan struct{})
	go func() {
		m.WorkerPool.Stop()
		close(done)
	}()

	// 최대 10초만 기다리고, 안 끝나면 강제 종료
	select {
	case <-done:
		log.Println("Cleanup finished successfully.")
	case <-time.After(10 * time.Second):
		log.Println("Cleanup timeout! Forced shutdown.")
	}
}
