package main

import (
	"context"
	"log"
	"message/internal/delivery/router"
	"message/internal/di"
	"message/internal/infrastructure/config"
	"message/internal/infrastructure/logger"
	"message/internal/infrastructure/migration"
	"message/internal/infrastructure/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// 1. 서버 및 모듈 초기화
	server, modules := InitServer()

	// 2. 서버 실행 (비동기)
	go func() {
		log.Println("Message service is running on :8083")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 3. 시스템 시그널 대기 (SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Message service ...")

	// 4. Graceful Shutdown 실행
	// HTTP 서버 먼저 종료 (새로운 요청 차단)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Message service shutdown:", err)
	}

	// 5. 비동기 워커 풀 종료 (남은 작업 처리)
	modules.ChatModule.Cleanup()
	// 필요하다면 다른 모듈의 Cleanup도 호출
	// modules.NoteModule.Cleanup()

	log.Println("Message service exiting")
}

// 모듈들을 묶어서 반환하기 위한 구조체
type AppModules struct {
	ChatModule *di.ChatModule
	// 다른 모듈들도 Cleanup이 필요하면 여기에 추가
}

func InitServer() (*http.Server, *AppModules) {

	// ---- Server Config Init -----
	sfg := config.NewServerConfig()

	// ---- DB Connect -----
	db := config.ConnectDatabase(sfg)

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

	//dataLoader.Register(loader.NewDeviceTokenInfoLoader(db, deviceStorage))

	// ---- Router Init -----
	// SetDefaultRoutes() 안에서 새로운 gin.Engine을 매번 생성하면 각기 다른 서버 인스턴스가 됩니다.
	// 이런 경우는 서버를 2개 띄우는 것과 같으므로 주의.
	router := router.NewMessageRouter("message", sfg.TokenConfig, logger)

	// ---- Domain Handler Init -----

	noteModule := di.InitNoteModule(db, mb)
	router.SetNoteRoutes(noteModule.Handler)

	lineKeyModule := di.InitLineKeyModule(db)
	router.SetLineKeyRoutes(lineKeyModule.Handler)

	chatFileModule := di.InitChatFileModule(db)

	chatModule := di.InitChatModule(db, mb, logger)
	router.SetChatRoutes(chatModule.Handler)

	otpModule := di.InitOtpModule(db, otpStorage)
	router.SetOtpRoutes(otpModule.Handler)

	chatRoomModule := di.InitChatRoomModule(db, chatRoomStorage, mb, logger)
	router.SetChatRoomRoutes(chatRoomModule.Handler)

	chatRoomFixedModule := di.InitChatRoomFixedModule(db)

	chatRoomTitleModule := di.InitChatRoomTitleModule(db)
	router.SetChatRoomTitleRoutes(chatRoomTitleModule.Handler)

	chatRoomConfigModule := di.InitChatRoomConfigModule(db)

	// chatService에도 chatRoom이 들어가지만, 다른 usecase의 조합으로 처리해야 할 수 있으므로 chat과 chatRoom을 분리.
	// usecase의 조합이지만 메인이 뭐냐? 라고 생각하고 작업하기
	chatServiceModule := di.InitChatServiceModule(chatModule.Usecase, lineKeyModule.Usecase, chatRoomModule.Usecase)
	router.SetChatServiceRoutes(chatServiceModule)

	chatRoomServiceModule := di.InitChatRoomServiceModule(chatRoomModule.Usecase, lineKeyModule.Usecase, chatModule.Usecase, chatRoomFixedModule.Usecase, chatRoomTitleModule.Usecase, chatRoomConfigModule.Usecase)
	router.SetChatRoomServiceRoutes(chatRoomServiceModule)

	chatLineServiceModule := di.InitChatLineServiceModule(chatModule.Usecase, chatRoomModule.Usecase, chatFileModule.Usecase)
	router.SetChatLineServiceRoutes(chatLineServiceModule)

	server := &http.Server{
		Addr:    ":8083",
		Handler: router.GetEngine(),
	}

	return server, &AppModules{
		ChatModule: chatModule,
	}
}
