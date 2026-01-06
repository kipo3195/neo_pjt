package main

import (
	"log"
	"net/http"
	natsBrocker "notificator/internal/delivery/adapter/nats"
	router "notificator/internal/delivery/router"
	"notificator/internal/di"
	"notificator/internal/infrastructure/config"
	"notificator/internal/infrastructure/migration"
	"notificator/internal/infrastructure/sender"
	"notificator/internal/infrastructure/storage"
)

func main() {
	server := InitServer()
	if server != nil {
		log.Println("Notificator service is running on :8082")
		log.Fatal(server.ListenAndServe())
	} else {
		log.Println("[ERROR] Notificator service is not available")
	}

}

func InitServer() *http.Server {

	// ---- Server Config Init -----
	sfg := config.NewServerConfig()

	// ---- DB Connect -----
	db := config.ConnectDatabase(sfg)

	// ---- DB Migration -----
	if sfg.AutoMigrate {
		migration.RunAll(db)
	}

	// ---- Storage Init -----
	chatRoomStorage := storage.NewChatRoomStorage()
	sendConnectionStorage := storage.NewSendConnectionStorage()

	// ---- Websocket sender Init
	chatDataSender := sender.NewChatDataSender()
	//noteDataSender : sender.NewNoteDataSender() TODO

	// ---- Data Loader -----

	// ---- Router Init -----
	router := router.NewNotificatorRouter("notificator", sfg.TokenConfig)

	// ---- Domain Handler Init -----
	chatRoomModule := di.InitChatRoomModule(db, chatRoomStorage)

	chatModule := di.InitChatModule(db, chatRoomStorage, sendConnectionStorage)

	noteModule := di.InitNoteModule(db)

	loginModule := di.InitLoginModule(db)

	socketSendModule := di.InitSocketSendModule(chatDataSender, sendConnectionStorage, chatRoomStorage)

	// ---- Service Handler Init ----
	notificatorServiceModule := di.InitNotificatorServiceModule(chatRoomModule.Usecase, socketSendModule.Usecase, loginModule.Usecase, sfg.WebsocketConnectionConfig)
	router.SetNotificatorServiceRoutes(notificatorServiceModule)

	// ---- Message Broker init ----
	mb := config.ConnectMessageBroker(sfg)
	conn := mb

	// ---- Message Broker Subscribe ----
	// 각 도메인별 핸들러 정의
	chatSub := natsBrocker.NewNatsChatSubscriber(conn, chatModule.Usecase, socketSendModule.Usecase)
	noteSub := natsBrocker.NewNatsNoteSubscriber(conn, noteModule.Usecase, socketSendModule.Usecase)
	chatRoomSub := natsBrocker.NewNatsChatRoomSubscriber(conn, chatModule.Usecase, noteModule.Usecase, socketSendModule.Usecase)

	// ---- NATS Subscribe ----
	// 도메인별 토픽만 구독
	chatSub.AddSubscribe("chat.broadcast")
	noteSub.AddSubscribe("note.broadcast")
	chatRoomSub.AddSubscribe("chat.room.broadcast")
	chatRoomSub.AddQueueSubscribe("create.chat.room")

	return &http.Server{
		Addr:    ":8082",
		Handler: router.R,
	}
}
