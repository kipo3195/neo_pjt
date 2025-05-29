package main

import (
	"admin/config"
	"admin/handlers"
	"admin/repositories"
	"admin/routes"
	"admin/usecases"
	"log"
	"net/http"
)

func main() {

	log.Println("Admin service is running on :8089")
	server := InitServer()
	log.Fatal(server.ListenAndServe())
}

func InitServer() *http.Server {

	sfg := config.NewServerConfig()
	db := config.ConnectDatabase(sfg)

	adminRepo := repositories.NewAdminRepository(db)
	adminUC := usecases.NewAdminUsecase(adminRepo)
	adminHandler := handlers.NewAdminHandler(adminUC)

	router := routes.SetupRoutes(adminHandler)

	return &http.Server{
		Addr:    ":8089",
		Handler: router,
	}

}
