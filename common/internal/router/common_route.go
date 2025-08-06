package router

import (
	appToken "common/internal/domains/appToken"
	appValidation "common/internal/domains/appValidation"
	"common/internal/domains/device"
	skin "common/internal/domains/skin"

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
	client.POST("/refresh", handlers.ClientHandler.AppTokenRefresh)

}

func SetSkinRoutes(parent *gin.RouterGroup, handlers *skin.SkinHandlers) {
	client := parent.Group("/client/v1/skin-img")
	client.GET("/", handlers.ClientHandler.GetSkinImage)

	server := parent.Group("/client/v1/skin-img")
	server.POST("/", handlers.ServerHandler.PutSkinImg)
}

func SetDeviceRoute(parent *gin.RouterGroup, handlers *device.DeviceHandlers) {
	server := parent.Group("/server/v1/device")
	server.POST("/init", handlers.ServerHandler.DeviceInit)
}
