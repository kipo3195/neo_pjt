package main

import (
	"core/config"
	"core/handlers"
	"core/repositories"
	"core/routes"
	"core/usecases"
	"log"
	"net/http"
)

func main() {
	log.Println("core service is running on :8085")
	server := InitServer()
	log.Fatal(server.ListenAndServe())
}

func InitServer() *http.Server {

	sfg := config.NewServerConfig()
	db := config.ConnectDatabase(sfg)

	coreRepo := repositories.NewCoreRepository(db)
	coreUC := usecases.NewCoreUsecase(coreRepo)
	coreHandler := handlers.NewCoreHandler(coreUC)

	router := routes.SetupRoutes(coreHandler)

	return &http.Server{
		Addr:    ":8085",
		Handler: router,
	}
}
