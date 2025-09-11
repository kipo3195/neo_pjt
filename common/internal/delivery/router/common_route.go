package router

import (
	"common/internal/delivery/handler"
	"common/internal/delivery/middleware"

	"github.com/gin-gonic/gin"
)

func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	return r, r.Group("/" + serviceName)
}

func SetAppValidationRoutes(parent *gin.RouterGroup, handler *handler.AppValidationHandler) {

	server := parent.Group("/server/v1/app-validation")

	server.GET("/validate", handler.GetAppValidation)

}

func SetAppTokenRoutes(parent *gin.RouterGroup, handler *handler.AppTokenHandler) {

	client := parent.Group("/client/v1/app-token")
	client.Use(middleware.AuthMiddleware()) // JWT 적용
	client.POST("/refresh", handler.AppTokenRefresh)

}

func SetSkinRoutes(parent *gin.RouterGroup, handler *handler.SkinHandler) {
	client := parent.Group("/client/v1/skin-img")
	client.Use(middleware.AuthMiddleware()) // JWT 적용
	client.GET("/", handler.GetSkinImage)

	server := parent.Group("/server/v1/skin-img")
	server.POST("/", handler.PutSkinImg)
}

// 이 API도 service 처리.
// appValidation 했던것 처럼  config, skin을 아우르는  service생성하고 api의 전달되는 데이터에 따라서 조회하여 response하도록 수정.
func SetConfigurationRoutes(parent *gin.RouterGroup, handlers *handler.ConfigurationHandler) {
	// client := parent.Group("/client/v1/config-hash")
	// client.GET("/config-hash", handlers.ClientHandler.GetConfigHash)
}

func SetUserRoutes(parent *gin.RouterGroup, handler *handler.UserHandler) {
	client := parent.Group("client/v1/user/register")
	client.POST("/", handler.UserRegister)
	client.GET("/challenge", handler.GetUserRegisterChallenge)

}

func SetInitAppValidtaionRoutes(parent *gin.RouterGroup, handler *handler.AppValidationHandler) {
	client := parent.Group("client/v1/app-validation")
	client.Use(middleware.AuthMiddleware())
	client.GET("/", handler.GetAppValidation)
}

func SetDeviceRoutes(parent *gin.RouterGroup, handler *handler.DeviceHandler) {
	client := parent.Group("server/v1/device-init")
	client.GET("/", handler.DeviceInit)
}
