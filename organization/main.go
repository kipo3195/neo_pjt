package main

import (
	"log"
	"net/http"
	"org/config"
	"org/handlers"
	"org/repositories"
	"org/routes"
	"org/usecases"
)

func main() {
	log.Println("org service is running on :8088")
	server := InitServer()
	log.Fatal(server.ListenAndServe())
}

func InitServer() *http.Server {

	sfg := config.NewServerConfig()
	db := config.ConnectDatabase(sfg)

	orgRepo := repositories.NewOrgRepository(db)
	orgUsecase := usecases.NewOrgUsecase(orgRepo)
	orgHandler := handlers.NewOrgHandler(sfg, orgUsecase)

	userRepo := repositories.NewUserRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo)
	userHandler := handlers.NewUserHandler(sfg, userUsecase)

	serverRepo := repositories.NewServerRepository(db)
	serverUsecase := usecases.NewServerUsecase(serverRepo)
	serverHandler := handlers.NewServerHandler(sfg, serverUsecase)

	router := routes.SetupRoutes(orgHandler, userHandler, serverHandler)

	return &http.Server{
		Addr:    ":8088",
		Handler: router,
	}
}
