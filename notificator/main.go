package main

import (
	"log"
	"net/http"
	router "notificator/internal/delivery/router"
	"notificator/internal/di"
	"notificator/internal/infrastructure/config"
)

func main() {
	server := InitServer()
	if server != nil {
		log.Println("Message service is running on :8087")
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

	// ---- Message Broker init ----
	mb := config.ConnectMessageBroker(sfg)

	// ---- DB Migration -----

	// ---- Storage Init -----

	// ---- Data Loader -----

	// ---- Router Init -----
	router := router.NewNotificatorRouter("notificator")

	// ---- Domain Handler Init -----
	chatModule := di.InitChatModule(db, mb)

	noteModule := di.InitNoteModule(db, mb)

	// ---- Service Handler Init ----
	notificatorServiceModule := di.InitNotificatorServiceModule(chatModule.Usecase, noteModule.Usecase)
	router.SetNotificatorServiceRoutes(notificatorServiceModule)

	return &http.Server{
		Addr:    ":8087",
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
