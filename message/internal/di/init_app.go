package di

import (
	"fmt"
	"log"
	"message/internal/adapter/http/router"
	"message/internal/app/loader"
	"message/internal/infrastructure/config"
	"message/internal/infrastructure/logger"
	"message/internal/infrastructure/pb"
	"message/internal/infrastructure/persistence/migration"
	"message/internal/infrastructure/storage"
	"net"
	"net/http"

	"google.golang.org/grpc"
	// ... 필요한 임포트들
)

type AppContainer struct {
	Server                 *http.Server
	Cleanup                func()
	DataLoader             *loader.DataLoader
	BatchMessageGrpcServer *grpc.Server
	BatchMessageListener   net.Listener
}

// InitApp: main.go에서 호출할 최종 조립 함수
func InitApp() (*AppContainer, error) {

	// ---- Server Config Init -----
	sfg := config.NewServerConfig()

	// ---- DB Connect -----
	db, err := config.ConnectDatabase(sfg)
	if err != nil {
		return nil, fmt.Errorf("db connection failed: %w", err)
	}

	// ---- Redis Connect -----
	cacheClient, err := config.ConnectCacheDataBase(sfg)
	if err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	// ---- gRPC Connect -----
	fileServiceGrpcClient, err := config.NewFileServiceProtocolBufferClient(sfg)
	if err != nil {
		return nil, fmt.Errorf("grpc client failed: %w", err)
	}

	// ---- gRPC Server Port Listener
	batchMessageListener, err := config.GetBatchMessageLis()
	if err != nil {
		return nil, fmt.Errorf("failed to listen gRPC port: %w", err)
	}

	// ---- Message Broker init ----
	mb, err := config.ConnectMessageBroker(sfg)
	if err != nil {
		return nil, fmt.Errorf("message broker failed: %w", err)
	}

	// ---- DB Migration -----
	if sfg.AutoMigrate {
		migration.RunAll(db)
	}

	// ---- gRPC Server -----
	batchMessageGrpcServer := grpc.NewServer()

	// ---- LOGGER Init ----
	logger := logger.NewSlogLogger()

	// ---- Storage Init -----
	otpStorage := storage.NewOtpStorage()
	chatRoomStorage := storage.NewChatRoomStorage()

	// ---- Data Loader -----
	dataLoader := loader.NewDataLoader()
	//dataLoader.Register(loader.NewDeviceTokenInfoLoader(db, deviceStorage)) 혹은 생성자 주입

	// ---- Domain Module Init -----
	noteModule := InitNoteModule(db, mb)
	lineKeyModule := InitLineKeyModule(db)
	chatFileModule := InitChatFileModule(db)
	chatModule := InitChatModule(db, mb, cacheClient, logger, fileServiceGrpcClient)
	otpModule := InitOtpModule(db, otpStorage)
	chatRoomModule := InitChatRoomModule(db, chatRoomStorage, mb, logger)
	chatRoomTitleModule := InitChatRoomTitleModule(db)
	chatRoomFixedModule := InitChatRoomFixedModule(db)
	chatRoomConfigModule := InitChatRoomConfigModule(db)

	// ---- Domain Service Module Init -----
	chatServiceModule := InitChatServiceModule(chatModule.Usecase, lineKeyModule.Usecase, chatRoomModule.Usecase)
	chatRoomServiceModule := InitChatRoomServiceModule(chatRoomModule.Usecase, lineKeyModule.Usecase, chatModule.Usecase, chatRoomFixedModule.Usecase, chatRoomTitleModule.Usecase, chatRoomConfigModule.Usecase)
	chatLineServiceModule := InitChatLineServiceModule(chatModule.Usecase, chatRoomModule.Usecase, chatFileModule.Usecase)

	// ---- Router Init -----
	// SetDefaultRoutes() 안에서 새로운 gin.Engine을 매번 생성하면 각기 다른 서버 인스턴스가 됩니다.
	// 이런 경우는 서버를 2개 띄우는 것과 같으므로 주의.
	router := router.NewMessageRouter("message", sfg.TokenConfig, logger)

	router.SetNoteRoutes(noteModule.Handler)
	router.SetLineKeyRoutes(lineKeyModule.Handler)
	router.SetChatRoutes(chatModule.Handler)
	router.SetOtpRoutes(otpModule.Handler)
	router.SetChatRoomRoutes(chatRoomModule.Handler)

	router.SetChatRoomTitleRoutes(chatRoomTitleModule.Handler)

	// chatService에도 chatRoom이 들어가지만, 다른 usecase의 조합으로 처리해야 할 수 있으므로 chat과 chatRoom을 분리.
	// usecase의 조합이지만 메인이 뭐냐? 라고 생각하고 작업하기
	router.SetChatServiceRoutes(chatServiceModule)
	router.SetChatRoomServiceRoutes(chatRoomServiceModule)
	router.SetChatLineServiceRoutes(chatLineServiceModule.Handler)

	// chatLineServiceModule는 하나의 서비스이지만 두개의 handler(http, gRPC)를 갖는다 -> 핵심 비즈니스 로직은 하나이며 분리될 수 없다.
	pb.RegisterBatchMessageServiceServer(batchMessageGrpcServer, chatLineServiceModule.GrpcHandler)

	// 자원 해제 - 실행 순서의 역순으로 종료 필요
	cleanup := func() {
		log.Println("--- Graceful Cleanup Start ---")

		if chatModule != nil {
			log.Println("Closing chatModule workers...")
			chatModule.Cleanup() // module에 있는 cleanup은 내부에 별도 고루틴을 통한 workerpool이 존재하는 경우 호출한다.
		}

		if mb != nil {
			log.Println("Closing Message Broker...")
			mb.Close()
		}

		if fileServiceGrpcClient != nil {
			log.Println("Closing gRPC client...")
			fileServiceGrpcClient.Close()
		}

		if cacheClient != nil {
			log.Println("Closing Redis client...")
			cacheClient.Close()
		}

		if db != nil {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				log.Println("Closing Database connection...")
				sqlDB.Close()
			}
		}

		log.Println("--- Graceful Cleanup Finished ---")
	}

	server := &http.Server{
		Addr:    ":8083",
		Handler: router.GetEngine(),
	}

	return &AppContainer{
		Server:                 server,
		Cleanup:                cleanup,
		DataLoader:             dataLoader,
		BatchMessageGrpcServer: batchMessageGrpcServer,
		BatchMessageListener:   batchMessageListener,
	}, nil
}
