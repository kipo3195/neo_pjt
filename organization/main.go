package main

import (
	"log"
	"net/http"
	"org/config"
	"org/handlers"
	"org/infra/storage"
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

	orgFileStorage := storage.NewOrgFileStorage() // 조직도 메모리 관리

	orgRepo := repositories.NewOrgRepository(db)
	orgUsecase := usecases.NewOrgUsecase(orgRepo, orgFileStorage)
	orgHandler := handlers.NewOrgHandler(sfg, orgUsecase)

	userRepo := repositories.NewUserRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo)
	userHandler := handlers.NewUserHandler(sfg, userUsecase)

	serverRepo := repositories.NewServerRepository(db)
	serverUsecase := usecases.NewServerUsecase(serverRepo, orgFileStorage)
	serverHandler := handlers.NewServerHandler(sfg, serverUsecase)

	handlers := &handlers.OrgHandlers{
		Org:    orgHandler,
		User:   userHandler,
		Server: serverHandler,
	}

	router := routes.SetupRoutes(handlers)

	return &http.Server{
		Addr:    ":8088",
		Handler: router,
	}
}
