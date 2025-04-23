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

	ruleRepo := repositories.NewRuleRepository(db)
	ruleUC := usecases.NewRuleUsecase(ruleRepo)
	ruleHandler := handlers.NewRuleHandler(ruleUC)

	authRepo := repositories.NewAuthRepository(db)
	authUC := usecases.NewAuthUsecase(authRepo)
	authHandler := handlers.NewAuthHandler(authUC)

	router := routes.SetupRoutes(ruleHandler, authHandler)

	return &http.Server{
		Addr:    ":8089",
		Handler: router,
	}

}
