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
	"log"
	"net/http"
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

	r, baseGroup := router.SetDefaultRoutes("common")

	appValidationHandler := appValidation.InitModule(db, configHashStorage)
	router.SetAppValidationRoutes(baseGroup, appValidationHandler)

	skinHandler := skin.InitModule(db, configHashStorage, skinStorage)
	router.SetSkinRoutes(baseGroup, skinHandler)

	appTokenHandler := appToken.InitModule(db)
	router.SetAppTokenRoutes(baseGroup, appTokenHandler)

	configurationHandler := configuration.InitModule(db, configHashStorage)
	router.SetConfigurationRoutes(baseGroup, configurationHandler)

	// service init
	deps := modules.Dependencies{
		DB:                db,
		ConfigHashStorage: configHashStorage,
	}

	appInitHandler := modules.InitAppInitModule(deps)
	r.POST("/server/v1/app-validation", appInitHandler.GetAppValidation)

	// 해야할 것
	// 실제로 appValidation에서는 api 요청을 받지 않아도 된다.  usecase만 구현되면 됨.
	//

	return &http.Server{
		Addr:    ":8086",
		Handler: r,
	}
}
