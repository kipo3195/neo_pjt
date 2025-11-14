package main

import (
	"log"
	"net/http"
	natsBrocker "notificator/internal/delivery/adapter/nats"
	router "notificator/internal/delivery/router"
	"notificator/internal/di"
	"notificator/internal/infrastructure/config"
	"notificator/internal/infrastructure/storage"
)

func main() {
	server := InitServer()
	if server != nil {
		log.Println("Message service is running on :8082")
		log.Fatal(server.ListenAndServe())
	} else {
		log.Println("[ERROR] Message service is not available")
	}

}

func InitServer() *http.Server {

	// ---- Server Config Init -----
	sfg := config.NewServerConfig()

	// ---- DB Connect -----
	db := config.ConnectDatabase(sfg)

	// ---- DB Migration -----

	// ---- Storage Init -----
	chatUserStorage := storage.NewChatUserStorage()

	// ---- Data Loader -----

	// ---- Router Init -----
	router := router.NewNotificatorRouter("notificator")

	// ---- Domain Handler Init -----
	chatModule := di.InitChatModule(db, chatUserStorage)

	noteModule := di.InitNoteModule(db)

	// ---- Service Handler Init ----
	notificatorServiceModule := di.InitNotificatorServiceModule(chatModule.Usecase, noteModule.Usecase)
	router.SetNotificatorServiceRoutes(notificatorServiceModule)

	// ---- Message Broker init ----
	mb := config.ConnectMessageBroker(sfg)

	// ---- Message Broker Subscribe ----
	// 이 로직도 어떻게 처리 안되나?
	conn := mb
	//defer conn.Close()
	sub := natsBrocker.NewNatsSubscriber(conn, chatModule.Usecase)

	// nats subscribe - 단순 for문 처리시 블로킹, 별도의 go 루틴으로 분리 필수.
	sub.StartSubscribe("chat.message")

	return &http.Server{
		Addr:    ":8082",
		Handler: router.R,
	}
	// if db != nil && mb != nil {

	// 	chatRepo := repository.NewChatRepository(db)
	// 	chatUC := usecase.NewChatUsecase(chatRepo, mb)

	// 	authRepo := repository.NewAuthRepository(db)
	// 	authUC := usecase.NewAuthUsecase(authRepo)

	// 	noteRepo := repository.NewNoteRepository(db)
	// 	noteUC := usecase.NewNoteUsecase(noteRepo)

	// 	messageHandler := handler.NewMessageHandler(chatUC, authUC, noteUC, mb)
	// 	router := usecase.SetupRoutes(messageHandler)

	// } else {
	// 	return nil
	// }
}
