package di

import (
	"log"
	"message/internal/adapter/http/handler"
	"message/internal/application/usecase"
	"message/internal/domain/logger"
	"message/internal/infrastructure/external/rpc"
	"message/internal/infrastructure/persistence/cacheStorage"
	"message/internal/infrastructure/persistence/repository"
	"message/internal/infrastructure/workerPool"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type ChatModule struct {
	Handler    *handler.ChatHandler
	Usecase    usecase.ChatUsecase
	WorkerPool workerPool.ChatWorkerPool
}

func InitChatModule(db *gorm.DB, connector *nats.Conn, cacheClient *redis.ClusterClient, logger logger.Logger, gRPCClient *grpc.ClientConn) *ChatModule {

	cacheStorage := cacheStorage.NewChatCache(cacheClient)
	repository := repository.NewChatRepository(db, cacheStorage)

	// 이 영역에서 구현체를 생성하고 인터페이스 타입으로 Usecase에 주입합니다 .
	workerPool := workerPool.NewChatWorkerPool(10, repository)
	workerPool.Init()
	apiRepository := rpc.NewGrpcChatApiRepositoryImpl(gRPCClient)
	usecase := usecase.NewChatUsecase(repository, connector, workerPool, logger, apiRepository)
	handler := handler.NewChatHandler(usecase)

	return &ChatModule{
		Handler:    handler,
		Usecase:    usecase,
		WorkerPool: workerPool,
	}
}

func (m *ChatModule) Cleanup() bool {
	log.Println("Cleaning up ChatModule...")

	ok := m.WorkerPool.Stop()

	if !ok {
		log.Println("WorkerPool not fully drained. Potential message loss.")
	}

	log.Println("Cleanup finished successfully.")
	return ok
}
