package main

import (
	"file/internal/delivery/router"
	"file/internal/di"
	"file/internal/infrastructure/config"
	"file/internal/infrastructure/logger"
	"file/internal/infrastructure/migration"
	"log"
	"net/http"
	// OCI SDK 필수 패키지
)

func main() {

	log.Println("file service is running on :8091")
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

	// ---- LOGGER Init ----
	logger := logger.NewSlogLogger()

	// ---- Storage Init -----

	// ---- Data Loader -----

	// ---- Router Init -----
	// SetDefaultRoutes() 안에서 새로운 gin.Engine을 매번 생성하면 각기 다른 서버 인스턴스가 됩니다.
	// 이런 경우는 서버를 2개 띄우는 것과 같으므로 주의.
	router := router.NewFileRouter("file", sfg.TokenConfig, logger)

	// ---- Domain Handler Init -----
	fileUrlModule := di.InitFileUrlModule(db, sfg.OracleStorageConfig, logger)
	router.SetFileUrlRoutes(fileUrlModule.Handler)

	return &http.Server{
		Addr:    ":8091",
		Handler: router.GetEngine(),
	}

}
