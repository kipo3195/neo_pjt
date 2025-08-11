package main

import (
	"common/internal/config"
	"common/internal/domains/appToken"
	appValidation "common/internal/domains/appValidation"
	"common/internal/domains/configuration"
	"common/internal/domains/skin"
	"common/internal/infra/storage"
	"common/internal/modules"
	"common/internal/router"
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

	// 메모리 저장소 생성 (빈 상태)
	configHashStorage := storage.NewConfigHashStorage()
	skinStorage := storage.NewSkinStorage()

	deps := modules.Dependencies{
		DB:                db,
		ConfigHashStorage: configHashStorage,
		SkinStorage:       skinStorage,
	}

	// 데이터 로딩

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := dataLoader.LoadAllData(ctx); err != nil {
		log.Fatal("Failed to load initial data:", err)
	}

	// api init

	r, baseGroup := router.SetDefaultRoutes("common")

	appValidation.InitModule(db, configHashStorage)
	//router.SetAppValidationRoutes(baseGroup, appValidationHandler)

	skinHandler := skin.InitModule(db, configHashStorage, skinStorage)
	router.SetSkinRoutes(baseGroup, skinHandler)

	appTokenHandler := appToken.InitModule(db)
	router.SetAppTokenRoutes(baseGroup, appTokenHandler)

	configurationHandler := configuration.InitModule(db, configHashStorage)
	router.SetConfigurationRoutes(baseGroup, configurationHandler)

	// service init
	appInitHandler := modules.InitAppInitModule(deps)
	r.POST("/server/v1/app-validation", appInitHandler.GetAppValidation)

	return &http.Server{
		Addr:    ":8086",
		Handler: r,
	}
}
