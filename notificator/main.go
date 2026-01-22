package main

import (
	"context"
	"log"
	"net/http"
	natsBrocker "notificator/internal/delivery/adapter/nats"
	router "notificator/internal/delivery/router"
	"notificator/internal/di"
	"notificator/internal/infrastructure/config"
	"notificator/internal/infrastructure/logger"
	"notificator/internal/infrastructure/migration"
	"notificator/internal/infrastructure/sender"
	"notificator/internal/infrastructure/storage"
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
		log.Println("Notificator service is running on :8082")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 3. 시스템 시그널 대기 (SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Notificator service ...")

	// 4. Graceful Shutdown 실행
	// HTTP 서버 먼저 종료 (새로운 요청 차단)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Notificator service shutdown:", err)
	}

	// 5. 비동기 워커 풀 종료 (남은 작업 처리)
	modules.ChatModule.Cleanup()
	// 필요하다면 다른 모듈의 Cleanup도 호출
	// modules.NoteModule.Cleanup()

	log.Println("Notificator service exiting")
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
	conn := mb

	// ---- LOGGER Init ----
	logger := logger.NewSlogLogger()

	// ---- Storage Init -----
	chatRoomStorage := storage.NewChatRoomStorage()
	sendConnectionStorage := storage.NewSendConnectionStorage()

	// ---- Websocket sender Init
	messageSender := sender.NewMessageSender(sendConnectionStorage)

	// ---- Data Loader -----

	// ---- Router Init -----
	router := router.NewNotificatorRouter("notificator", sfg.TokenConfig, logger)

	// ---- Domain Handler Init -----
	chatRoomModule := di.InitChatRoomModule(db, chatRoomStorage, sendConnectionStorage, conn, messageSender)
	chatModule := di.InitChatModule(db, chatRoomStorage, messageSender)
	noteModule := di.InitNoteModule(db)
	loginModule := di.InitLoginModule(db)
	socketSendModule := di.InitSocketSendModule(sendConnectionStorage)
	serviceUsersModule := di.InitServiceUsersModule(db)

	// ---- Service Handler Init ----
	notificatorServiceModule := di.InitNotificatorServiceModule(chatRoomModule.Usecase, socketSendModule.Usecase, loginModule.Usecase, sfg.WebsocketConnectionConfig)
	router.SetNotificatorServiceRoutes(notificatorServiceModule)

	// ---- Message Broker Subscribe ----
	// 각 도메인별 핸들러 정의
	chatSub := natsBrocker.NewNatsChatSubscriber(conn, chatModule.Usecase, messageSender)
	noteSub := natsBrocker.NewNatsNoteSubscriber(conn, noteModule.Usecase, socketSendModule.Usecase)
	chatRoomSub := natsBrocker.NewNatsChatRoomSubscriber(conn, chatRoomModule.Usecase, socketSendModule.Usecase)
	serviceUsersSub := natsBrocker.NewNatsServiceUsersSubscriber(conn, serviceUsersModule.Usecase)

	// ---- NATS Subscribe ----
	// 도메인별 토픽만 구독
	chatSub.AddSubscribe("chat.broadcast")
	chatSub.AddSubscribe("chat.count.broadcast")
	noteSub.AddSubscribe("note.broadcast")
	chatRoomSub.AddSubscribe("chat.room.broadcast")
	chatRoomSub.AddQueueSubscribe("create.chat.room")
	serviceUsersSub.AddQueueSubscribe("users.registered")

	server := &http.Server{
		Addr:    ":8082",
		Handler: router.R,
	}
	return server, &AppModules{
		ChatModule: chatModule,
	}
}
