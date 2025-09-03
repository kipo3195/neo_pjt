package main

import (
	"core/internal/delivery/router"
	"core/internal/di"
	"core/internal/infrastructure/config"
	"core/internal/infrastructure/migration"
	"core/internal/infrastructure/storage"
	"log"
	"net/http"
)

func main() {
	log.Println("core service is running on :8085")
	server := InitServer()
	log.Fatal(server.ListenAndServe())
}

func InitServer() *http.Server {
	// ---- Server Config Init -----
	sfg := config.NewServerConfig()

	// ---- DB Connect -----
	db := config.ConnectDatabase(sfg)

	// ---- DB Migration -----
	if sfg.AutoMigrate {
		migration.RunAll(db)
	}
	// ---- Storage Init -----
	serverInfoStorage := storage.NewServerInfoStorage()

	// ---- Data Loader -----

	// ---- Router Init -----
	r, baseGroup := router.SetDefaultRoutes("core")

	// ---- Domain Handler Init -----
	appValidationHandler := di.InitAppValidationHandler(db, sfg, serverInfoStorage)
	router.SetAppValidationRoute(baseGroup, appValidationHandler.Handler)

	// ---- Orchestrator Init -----

	// HTTP 서버 설정 및 반환
	return &http.Server{
		Addr:    ":8085",
		Handler: r,
	}
}
