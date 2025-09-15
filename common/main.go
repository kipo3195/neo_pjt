package main

import (
	"common/internal/delivery/router"
	"common/internal/di"
	"common/internal/infrastructure/config"
	"common/internal/infrastructure/loader"
	"common/internal/infrastructure/migration"
	"common/internal/infrastructure/storage"
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
	userStorage := storage.NewUserStorage()

	// deps := dependencies.Dependency{
	// 	DB:                db,
	// 	ConfigHashStorage: configHashStorage,
	// 	SkinStorage:       skinStorage,
	// 	AutoMigrate:       sfg.AutoMigrate,
	// }

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

	skinModule := di.InitSkinModule(db, skinStorage)
	router.SetSkinRoutes(baseGroup, skinModule.Handler)

	appTokenModule := di.InitAppTokenModule(db)
	router.SetAppTokenRoutes(baseGroup, appTokenModule.Handler)

	configurationModule := di.InitConfigurationModule(db, configHashStorage)
	router.SetConfigurationRoutes(baseGroup, configurationModule.Handler)

	userModule := di.InitUserModule(db, userStorage)
	router.SetUserRoutes(baseGroup, userModule.Handler)

	worksInfoModule := di.InitWorksInfoModule(db)

	// ---- Service Init -----
	appValidationHandler := di.InitAppValidationHandler(nil, skinModule.Usecase, configurationModule.Usecase)
	router.SetInitAppValidtaionRoutes(baseGroup, appValidationHandler)

	deviceInitHandler := di.InitDeviceInitHandler(worksInfoModule.Usecase, skinModule.Usecase, configurationModule.Usecase, appTokenModule.Usecase)
	router.SetDeviceRoutes(baseGroup, deviceInitHandler)

	return &http.Server{
		Addr:    ":8086",
		Handler: r,
	}
}
