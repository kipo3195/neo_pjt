package di

import (
	"log"
	"message/internal/adapter/http/router"
	"message/internal/app/loader"
	"message/internal/infrastructure/config"
	"message/internal/infrastructure/logger"
	"message/internal/infrastructure/persistence/migration"
	"message/internal/infrastructure/storage"
	"net/http"
	// ... 필요한 임포트들
)

type AppContainer struct {
	Server     *http.Server
	Cleanup    func()
	DataLoader *loader.DataLoader
}

// InitApp: main.go에서 호출할 최종 조립 함수
func InitApp() (*AppContainer, error) {

	// ---- Server Config Init -----
	sfg := config.NewServerConfig()

	// ---- DB Connect -----
	db := config.ConnectDatabase(sfg)

	// ---- Redis Connect -----
	cacheClient := config.ConnectCacheDataBase(sfg)

	// ---- gRPC Connect -----
	gRPCClient, err := config.NewProtocolBufferClient(sfg)
	if err != nil {
		log.Fatal(err)
	}

	// ---- DB Migration -----
	if sfg.AutoMigrate {
		migration.RunAll(db)
	}

	// ---- Message Broker init ----
	mb := config.ConnectMessageBroker(sfg)

	// ---- LOGGER Init ----
	logger := logger.NewSlogLogger()

	// ---- Storage Init -----
	otpStorage := storage.NewOtpStorage()
	chatRoomStorage := storage.NewChatRoomStorage()
	// ---- Data Loader -----
	dataLoader := loader.NewDataLoader()
	//dataLoader.Register(loader.NewDeviceTokenInfoLoader(db, deviceStorage)) 혹은 생성자 주입

	// ---- Router Init -----
	// SetDefaultRoutes() 안에서 새로운 gin.Engine을 매번 생성하면 각기 다른 서버 인스턴스가 됩니다.
	// 이런 경우는 서버를 2개 띄우는 것과 같으므로 주의.
	router := router.NewMessageRouter("message", sfg.TokenConfig, logger)

	// ---- Domain Handler Init -----

	noteModule := InitNoteModule(db, mb)
	router.SetNoteRoutes(noteModule.Handler)

	lineKeyModule := InitLineKeyModule(db)
	router.SetLineKeyRoutes(lineKeyModule.Handler)

	chatFileModule := InitChatFileModule(db)

	chatModule := InitChatModule(db, mb, cacheClient, logger, gRPCClient)
	router.SetChatRoutes(chatModule.Handler)

	otpModule := InitOtpModule(db, otpStorage)
	router.SetOtpRoutes(otpModule.Handler)

	chatRoomModule := InitChatRoomModule(db, chatRoomStorage, mb, logger)
	router.SetChatRoomRoutes(chatRoomModule.Handler)

	chatRoomFixedModule := InitChatRoomFixedModule(db)

	chatRoomTitleModule := InitChatRoomTitleModule(db)
	router.SetChatRoomTitleRoutes(chatRoomTitleModule.Handler)

	chatRoomConfigModule := InitChatRoomConfigModule(db)

	// chatService에도 chatRoom이 들어가지만, 다른 usecase의 조합으로 처리해야 할 수 있으므로 chat과 chatRoom을 분리.
	// usecase의 조합이지만 메인이 뭐냐? 라고 생각하고 작업하기
	chatServiceModule := InitChatServiceModule(chatModule.Usecase, lineKeyModule.Usecase, chatRoomModule.Usecase)
	router.SetChatServiceRoutes(chatServiceModule)

	chatRoomServiceModule := InitChatRoomServiceModule(chatRoomModule.Usecase, lineKeyModule.Usecase, chatModule.Usecase, chatRoomFixedModule.Usecase, chatRoomTitleModule.Usecase, chatRoomConfigModule.Usecase)
	router.SetChatRoomServiceRoutes(chatRoomServiceModule)

	chatLineServiceModule := InitChatLineServiceModule(chatModule.Usecase, chatRoomModule.Usecase, chatFileModule.Usecase)
	router.SetChatLineServiceRoutes(chatLineServiceModule)

	cleanup := func() {
		// 자원 해제 로직
	}

	server := &http.Server{
		Addr:    ":8083",
		Handler: router.GetEngine(),
	}

	return &AppContainer{
		Server:     server,
		Cleanup:    cleanup,
		DataLoader: dataLoader,
	}, nil
}
