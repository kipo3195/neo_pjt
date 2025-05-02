package main

import (
	"common/config"
	"common/handlers"
	"common/repositories"
	"common/routes"
	"common/usecases"
	"log"
	"net/http"
)

func main() {
	log.Println("common service is running on :8086")
	server := InitServer()
	log.Fatal(server.ListenAndServe())
}

func InitServer() *http.Server {

	sfg := config.NewServerConfig()
	db := config.ConnectDatabase(sfg)

	commonRepo := repositories.NewCommonRepository(db)
	commonUC := usecases.NewCommonUsecase(commonRepo)
	commonHandler := handlers.NewCommonHandler(commonUC)

	router := routes.SetupRoutes(commonHandler)

	return &http.Server{
		Addr:    ":8086",
		Handler: router,
	}
}
