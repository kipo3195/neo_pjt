package di

import (
	"file/internal/adapter/http/router"
	"file/internal/app/loader"
	"file/internal/infrastructure/config"
	"file/internal/infrastructure/logger"
	"file/internal/infrastructure/pb"
	"file/internal/infrastructure/persistence/migration"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

type AppContainer struct {
	Server                *http.Server
	Cleanup               func()
	DataLoader            *loader.DataLoader
	MessageFileGrpcServer *grpc.Server
	MessageFileListener   net.Listener
	BatchFileGrpcServer   *grpc.Server
	BatchFileListener     net.Listener
}

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

	// ---- gRPC Listener Init -----
	messageFileLis, err := config.GetMessageFileLis()
	log.Printf("message - file 실제 바인딩된 주소: %s", messageFileLis.Addr().String())
	if err != nil {
		return nil, fmt.Errorf("failed to listen gRPC port: %w", err)
	}

	batchFileLis, err := config.GetBatchFileLis()
	log.Printf("batch - file 실제 바인딩된 주소: %s", batchFileLis.Addr().String())
	if err != nil {
		return nil, fmt.Errorf("failed to listen gRPC port: %w", err)
	}

	// ---- DB Migration -----
	if sfg.AutoMigrate {
		migration.RunAll(db)
	}

	// ---- Message Broker init ----
	// mb, err := config.ConnectMessageBroker(sfg)
	// if err != nil {
	// 	return nil, fmt.Errorf("message broker failed: %w", err)
	// }

	// ---- DB Migration -----
	if sfg.AutoMigrate {
		migration.RunAll(db)
	}

	messageFileGrpcServer := grpc.NewServer(
	// 필요 시 인터셉터(Middleware) 추가 가능
	// grpc.UnaryInterceptor(authInterceptor),
	)

	batchFileGrpcServer := grpc.NewServer(
	// 필요 시 인터셉터(Middleware) 추가 가능
	// grpc.UnaryInterceptor(authInterceptor),
	)

	// ---- LOGGER Init ----
	logger := logger.NewSlogLogger()

	// ---- Storage Init -----

	// ---- Data Loader -----

	// ---- Router Init -----
	// SetDefaultRoutes() 안에서 새로운 gin.Engine을 매번 생성하면 각기 다른 서버 인스턴스가 됩니다.
	// 이런 경우는 서버를 2개 띄우는 것과 같으므로 주의.
	router := router.NewFileRouter("file", sfg.TokenConfig, logger)

	// ---- Domain Module Init -----
	fileUrlModule := InitFileUrlModule(db, cacheClient, sfg.OracleStorageConfig, logger)
	chatFileModule := InitChatFileModule(db)
	uploadFileCheckModule := InitUploadFileCheckModule(db)

	// ---- Domain Service Module Init -----

	uploadFileCheckServiceModule := InitUploadFileCheckServiceModule(chatFileModule.Usecase, uploadFileCheckModule.Usecase)

	// ---- Router Init -----
	router.SetFileUrlRoutes(fileUrlModule.Handler)

	// ---- gRPC Init 서비스를 다른걸로 띄워줘야함.. 필수 -----
	pb.RegisterFileServiceServer(messageFileGrpcServer, chatFileModule.ChatFileGrpcHandler)
	pb.RegisterUploadFileCheckServiceServer(batchFileGrpcServer, uploadFileCheckServiceModule)

	// 자원 해제 - 실행 순서의 역순으로 종료 필요
	cleanup := func() {
		log.Println("--- Graceful Cleanup Start ---")

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
		Addr:    ":8091",
		Handler: router.GetEngine(),
	}

	return &AppContainer{
		Server:  server,
		Cleanup: cleanup,
		//DataLoader: dataLoader,
		MessageFileGrpcServer: messageFileGrpcServer,
		MessageFileListener:   messageFileLis,
		BatchFileGrpcServer:   batchFileGrpcServer,
		BatchFileListener:     batchFileLis,
	}, nil

}
