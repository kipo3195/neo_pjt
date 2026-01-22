package main

import (
	"core/internal/delivery/router"
	"core/internal/di"
	"core/internal/infrastructure/config"
	"core/internal/infrastructure/logger"
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

	// ---- LOGGER Init ----
	logger := logger.NewSlogLogger()

	// ---- Data Loader -----

	// ---- Router Init -----
	router := router.NewCoreRouter("core", logger)
	// SetDefaultRoutes() 안에서 새로운 gin.Engine을 매번 생성하면 각기 다른 서버 인스턴스가 됩니다.
	// 이런 경우는 서버를 2개 띄우는 것과 같으므로 주의.

	// ---- Domain Handler Init -----
	appValidationHandler := di.InitAppValidationHandler(db, sfg, serverInfoStorage)
	router.SetAppValidationRoute(appValidationHandler.Handler)

	// ---- Orchestrator Init -----

	// HTTP 서버 설정 및 반환
	return &http.Server{
		Addr:    ":8085",
		Handler: router.GetEngine(),
	}
}
