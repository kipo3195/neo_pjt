package main

import (
	"auth/internal/handlers"
	"auth/internal/repositories"
	"auth/internal/routes"
	"auth/internal/usecases"
	"auth/pkg/config"
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

	baseGroup := routes.SetDefaultRoutes("auth")

	// 이런 구조로 변경할것. https://chatgpt.com/c/688325b2-b234-8005-856c-8ea21fde3fb3
	routes.SetClientCertificationRoute(baseGroup)
	routes.SetupServerTokenRoute(baseGroup)

	authRepo := repositories.NewAuthRepository(db)
	authUsecase := usecases.NewAuthUsecase(authRepo, sfg.GetJWTConfig())
	authHandler := handlers.NewAuthHandler(authUsecase)

	serverRepo := repositories.NewServerRepository(db)
	serverUsecase := usecases.NewServerUsecase(serverRepo, authRepo, sfg.GetJWTConfig())
	serverHandler := handlers.NewServerHandler(serverUsecase)

	router := routes.SetupRoutes(authHandler, serverHandler)

	return &http.Server{
		Addr:    ":8087",
		Handler: router,
	}
}
