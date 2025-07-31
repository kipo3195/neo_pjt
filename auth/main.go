package main

import (
	"auth/internal/config"
	certification "auth/internal/domains/certification"
	token "auth/internal/domains/token"
	"auth/internal/routes"
	"auth/internal/utils"
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

	r, baseGroup := routes.SetDefaultRoutes("auth")
	// SetDefaultRoutes() 안에서 새로운 gin.Engine을 매번 생성하면 각기 다른 서버 인스턴스가 됩니다.
	// 이런 경우는 서버를 2개 띄우는 것과 같으므로 주의.

	authUtil := utils.NewAuthUtil(sfg.GetJWTConfig())

	// 이런 구조로 변경할것.
	certificationHandler := certification.InitModules(db, sfg.GetJWTConfig(), authUtil)
	routes.SetLoginRoute(baseGroup, certificationHandler)

	tokenHandler := token.InitModules(db, sfg.GetJWTConfig(), authUtil)
	routes.SetTokenRoute(baseGroup, tokenHandler)

	return &http.Server{
		Addr:    ":8087",
		Handler: r,
	}
}
