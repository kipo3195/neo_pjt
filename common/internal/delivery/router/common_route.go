package router

import (
	"common/internal/delivery/handler"
	"common/internal/delivery/middleware"

	"github.com/gin-gonic/gin"
)

type commonRouter struct {
	R      *gin.Engine
	parent *gin.RouterGroup
}

type CommonRouter interface {
	SetAppValidationRoutes(handler *handler.AppValidationServiceHandler)
	SetAppTokenRoutes(handler *handler.AppTokenHandler)
	SetSkinRoutes(handler *handler.SkinHandler)
	SetConfigurationRoutes(handlers *handler.ConfigurationHandler)
	SetUserRoutes(handler *handler.UserHandler)
	SetInitAppValidtaionRoutes(handler *handler.AppValidationServiceHandler)
	SetDeviceRoutes(handler *handler.DeviceHandler)
	SetProfileRoutes(handler *handler.ProfileHandler)
	GetEngine() *gin.Engine
}

func NewCommonRouter(serviceName string) CommonRouter {
	r := gin.Default()
	parent := r.Group("/" + serviceName)
	return &commonRouter{
		R:      r,
		parent: parent,
	}
}

func (r *commonRouter) GetEngine() *gin.Engine {
	return r.R
}

func (r *commonRouter) SetAppValidationRoutes(handler *handler.AppValidationServiceHandler) {

	server := r.parent.Group("/server/v1/app-validation")

	server.GET("/validate", handler.GetAppValidation)

}

func (r *commonRouter) SetAppTokenRoutes(handler *handler.AppTokenHandler) {

	client := r.parent.Group("/client/v1/app-token")
	client.Use(middleware.AuthMiddleware()) // JWT 적용
	client.POST("/refresh", handler.AppTokenRefresh)

}

func (r *commonRouter) SetSkinRoutes(handler *handler.SkinHandler) {
	client := r.parent.Group("/client/v1/skin-img")
	client.Use(middleware.AuthMiddleware()) // JWT 적용
	client.GET("/", handler.GetSkinImage)

	server := r.parent.Group("/server/v1/skin-img")
	server.POST("/", handler.PutSkinImg)
}

// 이 API도 service 처리.
// appValidation 했던것 처럼  config, skin을 아우르는  service생성하고 api의 전달되는 데이터에 따라서 조회하여 response하도록 수정.
func (r *commonRouter) SetConfigurationRoutes(handlers *handler.ConfigurationHandler) {
	// client := parent.Group("/client/v1/config-hash")
	// client.GET("/config-hash", handlers.ClientHandler.GetConfigHash)
}

func (r *commonRouter) SetUserRoutes(handler *handler.UserHandler) {
	client := r.parent.Group("client/v1/user/register")
	client.POST("/", handler.UserRegister)
	client.GET("/challenge", handler.GetUserRegisterChallenge)

}

func (r *commonRouter) SetInitAppValidtaionRoutes(handler *handler.AppValidationServiceHandler) {
	client := r.parent.Group("client/v1/app-validation")
	client.Use(middleware.AuthMiddleware())
	client.GET("/", handler.GetAppValidation)
}

func (r *commonRouter) SetDeviceRoutes(handler *handler.DeviceHandler) {
	server := r.parent.Group("server/v1/device-init")
	server.POST("/", handler.DeviceInit)
}

func (r *commonRouter) SetProfileRoutes(handler *handler.ProfileHandler) {
	client := r.parent.Group("client/v1/profile")
	client.Use(middleware.AuthMiddleware())
	client.PUT("/img", handler.UploadProfileImg)
	client.DELETE("/img", handler.DeleteProfileImg) // 기본 이미지로 변경

	client.POST("/msg", handler.RegistProfileMsg)

}
