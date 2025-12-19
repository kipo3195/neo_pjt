package main

import (
	"log"
	"net/http"
	"user/internal/delivery/router"
	"user/internal/di"
	"user/internal/infrastructure/config"
	"user/internal/infrastructure/migration"
	"user/internal/infrastructure/storage"
)

func main() {
	log.Println("User service is running on :8084")
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
	profileCacheStorage := storage.NewProfileCacheStorage()
	profileStorage := storage.NewServerProfileStorage("") // 추후 s3 storage로 전환
	userInfoServiceStorage := storage.NewUserInfoServiceStorage()

	// ---- Data Loader -----

	//dataLoader.Register(loader.NewDeviceTokenInfoLoader(db, deviceStorage))

	// ---- Router Init -----
	router := router.NewUserRouter("user", sfg.TokenConfig)
	// SetDefaultRoutes() 안에서 새로운 gin.Engine을 매번 생성하면 각기 다른 서버 인스턴스가 됩니다.
	// 이런 경우는 서버를 2개 띄우는 것과 같으므로 주의.

	// ---- Domain Handler Init -----
	profileModule := di.InitProfileModule(db, profileStorage, profileCacheStorage)
	router.SetProfileRoutes(profileModule.Handler)

	userDetailModule := di.InitUserDetailModule(db, userInfoServiceStorage)
	router.SetUserDetailRoutes(userDetailModule.Handler)

	// ---- Service Handler Init ----
	userInfoServiceModule := di.InitUserInfoServiceModule(profileModule.Usecase, userDetailModule.Usecase)
	router.SetUserInfoServiceRoutes(userInfoServiceModule)

	userBatchServiceModule := di.InitUserBatchServiceModule(profileModule.Usecase, userDetailModule.Usecase)
	router.SetUserBatchServiceRoute(userBatchServiceModule)

	return &http.Server{
		Addr:    ":8084",
		Handler: router.GetEngine(),
	}
}
