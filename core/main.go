package main

import (
	"core/config"
	"core/handlers"
	"core/infra/storage"
	"core/repositories"
	"core/routes"
	"core/usecases"
	"log"
	"net/http"
)

func main() {
	log.Println("core service is running on :8085")
	server := InitServer()
	log.Fatal(server.ListenAndServe())
}

func InitServer() *http.Server {
	// 서버 설정 로드 (.env 또는 환경 변수 기반)
	sfg := config.NewServerConfig()

	// 데이터베이스 연결 설정
	db := config.ConnectDatabase(sfg)

	// 서버 관련 메타 정보 메모리 캐시 초기화, 필요시 메모리 로딩 로직 (loader) 추가
	serverInfoStorage := storage.NewServerInfoStorage()

	coreRepo := repositories.NewCoreRepository(db)
	coreUsecase := usecases.NewCoreUsecase(coreRepo, serverInfoStorage)
	coreHandler := handlers.NewCoreHandler(sfg, coreUsecase)

	// 핸들러 구조체 초기화 - 여러개의 세부 서비스 대응
	handlers := &handlers.CoreHandlers{
		Core: coreHandler,
	}

	// 라우터 설정
	router := routes.SetupRoutes(handlers)

	// HTTP 서버 설정 및 반환
	return &http.Server{
		Addr:    ":8085",
		Handler: router,
	}
}
