package main

import (
	"auth/internal/delivery/router"
	"auth/internal/di"
	"auth/internal/infrastructure/config"
	"auth/internal/infrastructure/loader"
	"auth/internal/infrastructure/migration"
	"auth/internal/infrastructure/storage"
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Auth service is running on :8087")
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
	serviceUserStorage := storage.NewServiceUserStorage()
	userAuthStorage := storage.NewUserAuthStorage()
	deviceStorage := storage.NewDeviceStorage()
	authTokenStorage := storage.NewAuthTokenStorage()

	// ---- Data Loader -----
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dataLoader := loader.NewDataLoader()
	dataLoader.Register(loader.NewAuthTokenLoader(db, authTokenStorage)) // 인증 토큰 정보 DB에서 업로드
	dataLoader.Register(loader.NewServiceUserLoader(db, serviceUserStorage))
	//dataLoader.Register(loader.NewDeviceTokenInfoLoader(db, deviceStorage))

	if err := dataLoader.LoadAllData(ctx); err != nil {
		log.Fatal("Failed to load initial data:", err)
	}

	// ---- Router Init -----
	router := router.NewAuthRouter("auth", sfg.TokenConfig)
	// SetDefaultRoutes() 안에서 새로운 gin.Engine을 매번 생성하면 각기 다른 서버 인스턴스가 됩니다.
	// 이런 경우는 서버를 2개 띄우는 것과 같으므로 주의.

	// ---- Domain Handler Init -----
	// 이런 구조로 변경할것.
	tokenModule := di.InitTokenModule(db, sfg, authTokenStorage, serviceUserStorage)
	router.SetTokenRoutes(tokenModule.Handler)

	userAuthModule := di.InitUserAuthModule(db, userAuthStorage)
	router.SetUserAuthRoutes(userAuthModule.Handler)

	deviceModule := di.InitDeviceModule(db, deviceStorage, sfg.TokenConfig.AccessTokenHash, sfg.TokenConfig.RefreshTokenHash)
	router.SetDeviceRoutes(deviceModule.Handler)

	// ---- Service Handler Init ----
	userAuthServiceModule := di.InitUserAuthServiceModule(userAuthModule.Usecase, deviceModule.Usecase, tokenModule.Usecase)
	router.SetUserAuthServiceRoutes(userAuthServiceModule)

	userDeviceAuthServiceModule := di.InitDeviceAuthServiceModule(tokenModule.Usecase, deviceModule.Usecase)
	router.SetUserAuthDeviceServiceRoutes(userDeviceAuthServiceModule)

	return &http.Server{
		Addr:    ":8087",
		Handler: router.GetEngine(),
	}
}
