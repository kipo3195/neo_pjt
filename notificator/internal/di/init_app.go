package di

import (
	"fmt"
	"net/http"
	"notificator/internal/adapter/http/router"
	"notificator/internal/adapter/nats/subscriber"
	"notificator/internal/infrastructure/config"
	"notificator/internal/infrastructure/logger"
	"notificator/internal/infrastructure/persistence/migration"
	"notificator/internal/infrastructure/persistence/storage"
	"notificator/internal/infrastructure/sender"
)

// 모듈들을 묶어서 반환하기 위한 구조체
type AppContainer struct {
	Server     *http.Server
	ChatModule *ChatModule
	// 다른 모듈들도 Cleanup이 필요하면 여기에 추가
}

func InitApp() (*AppContainer, error) {

	// ---- Server Config Init -----
	sfg := config.NewServerConfig()

	// ---- DB Connect -----
	db, err := config.ConnectDatabase(sfg)
	if err != nil {
		return nil, fmt.Errorf("db connection failed: %w", err)
	}

	// ---- Message Broker init ----
	mb, err := config.ConnectMessageBroker(sfg)
	if err != nil {
		return nil, fmt.Errorf("message broker connect failed: %w", err)
	}

	// ---- LOGGER Init ----
	logger := logger.NewSlogLogger()

	// ---- DB Migration -----
	if sfg.AutoMigrate {
		migration.RunAll(db)
	}

	// ---- Storage Init -----
	chatRoomStorage := storage.NewChatRoomStorage()
	sendConnectionStorage := storage.NewSendConnectionStorage()

	// ---- Websocket sender Init
	messageSender := sender.NewMessageSender(sendConnectionStorage)

	// ---- Data Loader -----

	// ---- Router Init -----
	router := router.NewNotificatorRouter("notificator", sfg.TokenConfig, logger)

	// ---- Domain Handler Init -----
	chatRoomModule := InitChatRoomModule(db, chatRoomStorage, sendConnectionStorage, mb, messageSender)
	chatModule := InitChatModule(db, chatRoomStorage, messageSender)
	noteModule := InitNoteModule(db)
	loginModule := InitLoginModule(db)
	socketSendModule := InitSocketSendModule(sendConnectionStorage)
	serviceUsersModule := InitServiceUsersModule(db)

	// ---- Service Handler Init ----
	notificatorServiceModule := InitNotificatorServiceModule(chatRoomModule.Usecase, socketSendModule.Usecase, loginModule.Usecase, sfg.WebsocketConnectionConfig)

	router.SetNotificatorServiceRoutes(notificatorServiceModule)

	// ---- Message Broker Subscribe ----
	// 각 도메인별 핸들러 정의
	chatSub := subscriber.NewNatsChatSubscriber(mb, chatModule.Usecase, messageSender)
	noteSub := subscriber.NewNatsNoteSubscriber(mb, noteModule.Usecase, socketSendModule.Usecase)
	chatRoomSub := subscriber.NewNatsChatRoomSubscriber(mb, chatRoomModule.Usecase, socketSendModule.Usecase)
	serviceUsersSub := subscriber.NewNatsServiceUsersSubscriber(mb, serviceUsersModule.Usecase)

	// ---- NATS Subscribe ----
	// 도메인별 토픽만 구독
	chatSub.AddSubscribe("chat.broadcast")
	chatSub.AddSubscribe("chat.count.broadcast")
	chatSub.AddSubscribe("chat.read.broadcast")
	noteSub.AddSubscribe("note.broadcast")
	chatRoomSub.AddSubscribe("chat.room.broadcast")
	chatRoomSub.AddQueueSubscribe("create.chat.room")
	serviceUsersSub.AddQueueSubscribe("users.registered")

	server := &http.Server{
		Addr:    ":8082",
		Handler: router.R,
	}

	return &AppContainer{
		ChatModule: chatModule,
		Server:     server,
	}, nil
}
