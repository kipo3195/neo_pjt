package main

import (
	"common/handlers"
	loader "common/infra/loader"
	"common/internal/config"
	"common/internal/domains/appToken"
	appValidation "common/internal/domains/appValidation"
	"common/internal/domains/skin"
	"common/internal/infra/storage"
	"common/internal/router"
	"common/repositories"
	"common/usecases"
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

	publicRepo := repositories.NewPublicRepository(db)
	publicUsecase := usecases.NewPublicUsecase(publicRepo, configHashStorage)
	pubHandler := handlers.NewPublicHandler(publicUsecase)

	serverRepo := repositories.NewServerRepository(db)
	serverUsecase := usecases.NewServerUsecase(serverRepo, configHashStorage)
	serverHandler := handlers.NewServerHandler(serverUsecase)

	// DB 의존성 구성
	commonRepo := repositories.NewCommonRepository(db)
	// 의존성 주입 완료된 usecase 생성
	commonUsecase := usecases.NewCommonUsecase(commonRepo, configHashStorage)
	// usecase 내부 초기화 실행 (ex. DB → 캐시 로딩)
	commonLoader := loader.NewCommonLoader(commonUsecase)
	if err := commonLoader.RunAll(); err != nil {
		log.Fatalf("common loader 초기화 실패: %v", err)
		// 서버 종료됨.
	}
	log.Println("common service memory loading success !")
	// 초기화 완료된 usecase를 주입해 안전한 handler 구성
	commonHandler := handlers.NewCommonHandler(commonUsecase)

	// handlers := &handlers.CommonHandlers{
	// 	Common: commonHandler,
	// 	Server: serverHandler,
	// 	Public: pubHandler,
	// }

	// router := routes.SetupRoutes(handlers)

	return &http.Server{
		Addr:    ":8086",
		Handler: r,
	}
}
