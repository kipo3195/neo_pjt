package main

import (
	"log"
	"net/http"
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

	// ---- Domain Handler Init -----

	// ---- Service Handler Init ----

	if db != nil && mb != nil {

		chatRepo := repositories.NewChatRepository(db)
		chatUC := usecases.NewChatUsecase(chatRepo, mb)

		authRepo := repositories.NewAuthRepository(db)
		authUC := usecases.NewAuthUsecase(authRepo)

		noteRepo := repositories.NewNoteRepository(db)
		noteUC := usecases.NewNoteUsecase(noteRepo)

		messageHandler := handlers.NewMessageHandler(chatUC, authUC, noteUC, mb)
		router := routes.SetupRoutes(messageHandler)

		return &http.Server{
			Addr:    ":8087",
			Handler: router,
		}
	} else {
		return nil
	}
}
