package main

import (
	"core/internal/config"
	appValidation "core/internal/domains/appValidation"
	"core/internal/infra/storage"
	"core/internal/router"
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

	r, baseGroup := router.SetDefaultRoutes("core")

	appValidationHandlers := appValidation.InitModules(db, sfg, serverInfoStorage)
	router.SetAppValidationRoutes(baseGroup, appValidationHandlers)

	// HTTP 서버 설정 및 반환
	return &http.Server{
		Addr:    ":8085",
		Handler: r,
	}
}
