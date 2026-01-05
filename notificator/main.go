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
	chatUserStorage := storage.NewChatUserStorage()
	noteUserStorage := storage.NewNoteUserStorage()
	sendConnectionStorage := storage.NewSendConnectionStorage()

	// ---- Websocket sender Init
	chatDataSender := sender.NewChatDataSender()
	//noteDataSender : sender.NewNoteDataSender() TODO

	// ---- Data Loader -----

	// ---- Router Init -----
	router := router.NewNotificatorRouter("notificator", sfg.TokenConfig)

	// ---- Domain Handler Init -----
	chatModule := di.InitChatModule(db, chatUserStorage, sendConnectionStorage)

	noteModule := di.InitNoteModule(db, noteUserStorage)

	loginModule := di.InitLoginModule(db)

	socketSendModule := di.InitSocketSendModule(chatDataSender, sendConnectionStorage, chatUserStorage)

	// ---- Service Handler Init ----
	notificatorServiceModule := di.InitNotificatorServiceModule(chatModule.Usecase, noteModule.Usecase, socketSendModule.Usecase, loginModule.Usecase, sfg.WebsocketConnectionConfig)
	router.SetNotificatorServiceRoutes(notificatorServiceModule)

	// ---- Message Broker init ----
	mb := config.ConnectMessageBroker(sfg)

	// ---- Message Broker Subscribe ----
	conn := mb
	sub := natsBrocker.NewNatsSubscriber(conn, chatModule.Usecase, noteModule.Usecase, socketSendModule.Usecase)

	sub.AddSubscribe("chat.message")
	sub.AddSubscribe("note.message")
	sub.AddQueueSubscribe("create.chat.room.message")

	return &http.Server{
		Addr:    ":8082",
		Handler: router.R,
	}
}
