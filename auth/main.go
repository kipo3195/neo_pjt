package main

import (
	"auth/config"
	"auth/handlers"
	"auth/repositories"
	"auth/routes"
	"auth/usecases"
	"log"
	"net/http"
)

func main() {
	log.Println("Auth service is running on :8087")
	server := InitServer()
	log.Fatal(server.ListenAndServe())
}

func InitServer() *http.Server {

	sfg := config.NewServerConfig()
	db := config.ConnectDatabase(sfg)

	authRepo := repositories.NewAuthRepository(db)
	authUsecase := usecases.NewAuthUsecase(authRepo, sfg.GetJWTConfig())
	authHandler := handlers.NewAuthHandler(authUsecase)

	serverRepo := repositories.NewServerRepository(db)
	serverUsecase := usecases.NewServerUsecase(serverRepo, authRepo)
	serverHandler := handlers.NewServerHandler(serverUsecase)

	router := routes.SetupRoutes(authHandler, serverHandler)

	return &http.Server{
		Addr:    ":8087",
		Handler: router,
	}
}
