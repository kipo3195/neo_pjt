package main

import (
	"common/internal/delivery/middleware"
	"common/internal/delivery/router"
	"common/internal/di"
	"common/internal/infrastructure/config"
	"common/internal/infrastructure/loader"
	"common/internal/infrastructure/migration"
	"common/internal/infrastructure/storage"
	"common/internal/services/dependencies"
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("common service is running on :8086")
	server := InitServer()
	log.Fatal(server.ListenAndServe())
}

func InitServer() *http.Server {

	sfg := config.NewServerConfig()
	db := config.ConnectDatabase(sfg)

	// ---- DB Migration -----
	if sfg.AutoMigrate {
		migration.RunAll(db)
	}

	// ---- Storage Init -----
	configHashStorage := storage.NewConfigHashStorage()
	skinStorage := storage.NewSkinStorage()

	deps := dependencies.Dependency{
		DB:                db,
		ConfigHashStorage: configHashStorage,
		SkinStorage:       skinStorage,
		AutoMigrate:       sfg.AutoMigrate,
	}

	// ---- Data Loader -----
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dataLoader := loader.NewDataLoader()
	dataLoader.Register(loader.NewSkinLoader(db, skinStorage))
	dataLoader.Register(loader.NewConfigHashLoader(db, configHashStorage))

	if err := dataLoader.LoadAllData(ctx); err != nil {
		log.Fatal("Failed to load initial data:", err)
	}

	// ---- Router Init -----

	r, baseGroup := router.SetDefaultRoutes("common")

	skinHandler := di.InitSkinHandler(db, configHashStorage, skinStorage)
	router.SetSkinRoutes(baseGroup, skinHandler.Handler)

	appTokenHandler := di.InitAppTokenHandler(db)
	router.SetAppTokenRoutes(baseGroup, appTokenHandler.Handler)

	configurationHandler := di.InitConfigurationHandler(db, configHashStorage)
	router.SetConfigurationRoutes(baseGroup, configurationHandler.Handler)

	// ---- Service Init -----
	appInitHandler := di.InitAppValidationService(deps)
	r.POST("/client/v1/app-validation",
		middleware.AuthMiddleware(),     // <- 여기서 JWT 미들웨어 적용
		appInitHandler.GetAppValidation, // 실제 서비스 핸들러
	)

	deviceInitHandler := di.InitDeviceInitService((deps))
	r.POST("/server/v1/device-init", deviceInitHandler.DeviceInit)

	return &http.Server{
		Addr:    ":8086",
		Handler: r,
	}
}
