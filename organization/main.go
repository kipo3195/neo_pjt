package main

import (
	"context"
	"log"
	"net/http"
	"org/internal/delivery/router"
	"org/internal/di"
	"org/internal/infrastructure/config"
	"org/internal/infrastructure/loader"
	"org/internal/infrastructure/logger"
	"org/internal/infrastructure/migration"
	"org/internal/infrastructure/storage"
	"time"
)

func main() {
	log.Println("org service is running on :8088")
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
	orgFileStorage := storage.NewOrgFileStorage() // 조직도 메모리 관리
	orgStorage := storage.NewOrgStorage()

	// ---- Data Loader -----
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dataLoader := loader.NewDataLoader()
	dataLoader.Register(loader.NewOrgLoader(db, orgStorage))

	if err := dataLoader.LoadAllData(ctx); err != nil {
		log.Fatal("Failed to load initial data:", err)
	}

	// ---- Router Init -----
	router := router.NewOrgRoute("org", sfg.TokenConfig, logger)
	// SetDefaultRoutes() 안에서 새로운 gin.Engine을 매번 생성하면 각기 다른 서버 인스턴스가 됩니다.
	// 이런 경우는 서버를 2개 띄우는 것과 같으므로 주의.

	// ---- Domain Handler Init -----
	departmentModule := di.InitDepartmentModule(db)
	router.SetDepartmentRoutes(departmentModule.Handler)

	orgModule := di.InitOrgModule(db, orgFileStorage, orgStorage)
	router.SetOrgRoute(orgModule.Handler)

	userModule := di.InitUserModule(db)
	router.SetUserRoute(userModule.Handler)

	// ---- Orchestrator Init -----
	dummyDataInitServiceModule := di.InitDummyDataServiceModule(departmentModule.Usecase, orgModule.Usecase, userModule.Usecase)
	router.SetDummyDataServiceRoute(dummyDataInitServiceModule)

	orgUserServiceModule := di.InitOrgUserServiceModule(orgModule.Usecase, userModule.Usecase)
	router.SetOrgUserServiceRoute(orgUserServiceModule)

	orgBatchServiceModule := di.InitOrgBatchServiceModule(departmentModule.Usecase, orgModule.Usecase, userModule.Usecase)
	router.SetOrgBatchServiceRoute(orgBatchServiceModule)

	return &http.Server{
		Addr:    ":8088",
		Handler: router.GetEngine(),
	}
}
