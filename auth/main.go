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
	authUC := usecases.NewAuthUsecase(authRepo, sfg.GetJWTConfig())
	authHandler := handlers.NewAuthHandler(authUC)

	router := routes.SetupRoutes(authHandler)

	return &http.Server{
		Addr:    ":8087",
		Handler: router,
	}
}
