package router

import (
	"auth/internal/delivery/handler"

	"github.com/gin-gonic/gin"
)

func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	return r, r.Group("/" + serviceName)
}

func SetCertificationRoutes(parent *gin.RouterGroup, handler *handler.CertificationHandler) {
	client := parent.Group("/client/v1/certification")
	client.POST("/login", handler.Login)

}

func SetTokenRoutes(parent *gin.RouterGroup, handler *handler.TokenHandler) {
	// 클라이언트 영역이 있다면
	// client := parent.Group("/client/v1/token")

	server := parent.Group("/server/v1/token")
	server.POST("/generate-app-token", handler.GenerateAppToken)
	server.POST("/app-token-validation", handler.AppTokenValidation)
	server.POST("/app-token-refresh", handler.AppTokenRefresh)

}

func SetUserAuthRoutes(parent *gin.RouterGroup, handler *handler.UserAuthHandler) {
	client := parent.Group("/client/v1/user/auth")
	client.GET("/challenge", handler.GenerateAuthChallenge)
	client.GET("/", handler.GetAuthStatus)
}
