package main

import (
	"common/internal/config"
	"common/internal/domains/appToken"
	"common/internal/domains/configuration"
	"common/internal/domains/skin"
	"common/internal/infra/loader"
	"common/internal/infra/migration"
	"common/internal/infra/storage"
	"common/internal/middleware"
	"common/internal/router"
	"common/internal/services/dependencies"
	"common/internal/services/serviceModules"
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

	// ---- Router Init-----

	r, baseGroup := router.SetDefaultRoutes("common")

	skinHandler := skin.InitModule(db, configHashStorage, skinStorage)
	router.SetSkinRoutes(baseGroup, skinHandler)

	appTokenHandler := appToken.InitModule(db)
	router.SetAppTokenRoutes(baseGroup, appTokenHandler)

	configurationHandler := configuration.InitModule(db, configHashStorage)
	router.SetConfigurationRoutes(baseGroup, configurationHandler)

	// ---- Service Init ----
	appInitHandler := serviceModules.InitAppValidationModule(deps)
	r.POST("/client/v1/app-validation",
		middleware.AuthMiddleware(),     // <- 여기서 JWT 미들웨어 적용
		appInitHandler.GetAppValidation, // 실제 서비스 핸들러
	)

	deviceInitHandler := serviceModules.InitDeviceInitModule((deps))
	r.POST("/server/v1/device-init", deviceInitHandler.DeviceInit)

	return &http.Server{
		Addr:    ":8086",
		Handler: r,
	}
}
