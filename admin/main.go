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
	orgRepo := repositories.NewOrgRepository(db)
	orgUC := usecases.NewOrgUsecase(orgRepo)
	orgHandler := handlers.NewOrgHandler(orgUC)

	// common - 스킨, 앱 설정 관리
	commonRepo := repositories.NewCommonRepository(db)
	commonUC := usecases.NewCommonUsecase(commonRepo)
	commonHandler := handlers.NewCommonHandler(commonUC)

	handlers := &handlers.AdminHandlers{
		Org:    orgHandler,
		Common: commonHandler,
	}

	router := routes.SetupRoutes(handlers)

	return &http.Server{
		Addr:    ":8089",
		Handler: router,
	}

}
