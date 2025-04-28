package main

import (
	"log"
	"message/config"
	"message/handlers"
	"message/repositories"
	"message/routes"
	"message/usecases"
	"net/http"
)

func main() {

	log.Println("Message service is running on :8087")
	server := InitServer()
	log.Fatal(server.ListenAndServe())

}

func InitServer() *http.Server {

	sfg := config.NewServerConfig()
	db := config.ConnectDatabase(sfg)

	messageRepo := repositories.NewChatRepository(db)
	messageUC := usecases.NewChatUsecase(messageRepo)
	messageHandler := handlers.NewChatHandler(messageUC)

	router := routes.SetupRoutes(messageHandler)

	return &http.Server{
		Addr:    ":8087",
		Handler: router,
	}
}
