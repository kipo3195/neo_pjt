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
	server := InitServer()
	if server != nil {
		log.Println("Message service is running on :8087")
		log.Fatal(server.ListenAndServe())
	} else {
		log.Println("[ERROR] Message service is not available")
	}

}

func InitServer() *http.Server {

	// 서버 설정 읽기
	sfg := config.NewServerConfig()
	// DB nil일 경우 처리 필요
	db := config.ConnectDatabase(sfg)
	// 메시지 브로커, nil일 경우 처리 필요.
	mb := config.ConnectMessageBroker(sfg)

	// TODO
	// authRepo, authUC
	// noteRepo, noteUC
	// messageHandler에 주입.

	if db != nil && mb != nil {

		chatRepo := repositories.NewChatRepository(db)
		chatUC := usecases.NewChatUsecase(chatRepo)

		authRepo := repositories.NewAuthRepository(db)
		authUC := usecases.NewAuthUsecase(authRepo)

		messageHandler := handlers.NewMessageHandler(chatUC, authUC, mb)
		router := routes.SetupRoutes(messageHandler)

		return &http.Server{
			Addr:    ":8087",
			Handler: router,
		}
	} else {
		return nil
	}
}
