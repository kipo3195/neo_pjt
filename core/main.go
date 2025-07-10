package main

import (
	"core/config"
	"core/handlers"
	"core/infra/storage"
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

	serverInfoStorage := storage.NewServerInfoStorage()

	coreRepo := repositories.NewCoreRepository(db)
	coreUsecase := usecases.NewCoreUsecase(coreRepo, serverInfoStorage)
	coreHandler := handlers.NewCoreHandler(sfg, coreUsecase)

	log.Println("core service memory loading success !")
	handlers := &handlers.CoreHandlers{
		Core: coreHandler,
	}

	router := routes.SetupRoutes(handlers)

	return &http.Server{
		Addr:    ":8085",
		Handler: router,
	}
}
