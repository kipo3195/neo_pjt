package main

import (
	certification "auth/internal/domains/certification"
	"auth/internal/handlers"
	"auth/internal/repositories"
	"auth/internal/routes"
	"auth/internal/usecases"
	"auth/pkg/config"
	"log"
	"net/http"
)

func main() {
	log.Println("Auth service is running on :8087")
	server := InitServer()
	log.Fatal(server.ListenAndServe())
}

func InitServer() *http.Server {

	sfg := config.NewServerConfig()
	db := config.ConnectDatabase(sfg)

	baseGroup := routes.SetDefaultRoutes("auth")
	// SetDefaultRoutes() 안에서 새로운 gin.Engine을 매번 생성하면 각기 다른 서버 인스턴스가 됩니다.
	// 이런 경우는 서버를 2개 띄우는 것과 같으므로 주의.

	// 이런 구조로 변경할것.
	certificationHandler := certification.InitCertificationModule(db, sfg.GetJWTConfig())
	routes.SetLoginRoute(baseGroup, certificationHandler)

	routes.SetupServerTokenRoute(baseGroup)
	authRepo := repositories.NewAuthRepository(db)
	authUsecase := usecases.NewAuthUsecase(authRepo, sfg.GetJWTConfig())
	authHandler := handlers.NewAuthHandler(authUsecase)

	serverRepo := repositories.NewServerRepository(db)
	serverUsecase := usecases.NewServerUsecase(serverRepo, authRepo, sfg.GetJWTConfig())
	serverHandler := handlers.NewServerHandler(serverUsecase)

	router := routes.SetupRoutes(authHandler, serverHandler)

	return &http.Server{
		Addr:    ":8087",
		Handler: router,
	}
}
