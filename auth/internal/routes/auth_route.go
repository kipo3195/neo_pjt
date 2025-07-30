package routes

import (
	"auth/internal/domains/certification"
	"auth/internal/domains/token"

	"github.com/gin-gonic/gin"
)

func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	return r, r.Group("/" + serviceName)
}

func SetLoginRoute(parent *gin.RouterGroup, handlers *certification.CertificationHandlers) {

	client := parent.Group("/client/v1/certification")
	client.POST("/login", handlers.ClientHandler.Login)

	// 서버 영역이 있다면
	// server := parent.Group("/server/v1/certification")

}

func SetTokenRoute(parent *gin.RouterGroup, handlers *token.TokenHandlers) {

	// 클라이언트 영역이 있다면
	// client := parent.Group("/client/v1/token")

	server := parent.Group("/server/v1/token")
	server.POST("/generate-app-token", handlers.ServerHandler.GenerateAppToken)
	server.POST("/app-token-validation", handlers.ServerHandler.AppTokenValidation)
	server.POST("/app-token-refresh", handlers.ServerHandler.AppTokenRefresh)

}
