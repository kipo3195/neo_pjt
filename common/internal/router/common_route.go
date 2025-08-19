package router

import (
	"common/internal/domains/appToken"
	"common/internal/domains/appValidation"
	"common/internal/domains/configuration"
	"common/internal/domains/skin"
	"common/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	return r, r.Group("/" + serviceName)
}

func SetAppValidationRoutes(parent *gin.RouterGroup, handlers *appValidation.AppValidationHandlers) {

	server := parent.Group("/server/v1/app-validation")

	server.GET("/validate", handlers.ServerHandler.AppValidation)

}

func SetAppTokenRoutes(parent *gin.RouterGroup, handlers *appToken.AppTokenHandlers) {

	client := parent.Group("/client/v1/app-token")
	client.Use(middleware.AuthMiddleware()) // JWT 적용
	client.POST("/refresh", handlers.ClientHandler.AppTokenRefresh)

}

func SetSkinRoutes(parent *gin.RouterGroup, handlers *skin.SkinHandlers) {
	client := parent.Group("/client/v1/skin-img")
	client.Use(middleware.AuthMiddleware()) // JWT 적용
	client.GET("/", handlers.ClientHandler.GetSkinImage)

	server := parent.Group("/server/v1/skin-img")
	server.POST("/", handlers.ServerHandler.PutSkinImg)
}

// 이 API도 service 처리.
// appValidation 했던것 처럼  config, skin을 아우르는  service생성하고 api의 전달되는 데이터에 따라서 조회하여 response하도록 수정.
func SetConfigurationRoutes(parent *gin.RouterGroup, handlers *configuration.ConfigurationHandlers) {
	// client := parent.Group("/client/v1/config-hash")
	// client.GET("/config-hash", handlers.ClientHandler.GetConfigHash)
}
