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

	// org - 조직도 관리
	adminOrgRepo := repositories.NewAdminOrgRepository(db)
	adminOrgUC := usecases.NewAdminOrgUsecase(adminOrgRepo)
	adminOrgHandler := handlers.NewAdminHandler(adminOrgUC)

	handlers := &handlers.AdminHandlers{
		AdminOrgHandler: adminOrgHandler,
	}

	router := routes.SetupRoutes(handlers)

	return &http.Server{
		Addr:    ":8089",
		Handler: router,
	}

}
