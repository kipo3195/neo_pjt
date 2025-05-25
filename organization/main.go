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
	orgUC := usecases.NewOrgUsecase(orgRepo)
	orgHandler := handlers.NewOrgHandler(sfg, orgUC)

	router := routes.SetupRoutes(orgHandler)

	return &http.Server{
		Addr:    ":8088",
		Handler: router,
	}
}
