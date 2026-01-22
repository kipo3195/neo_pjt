package router

import (
	"common/internal/delivery/handler"
	"common/internal/delivery/middleware"
	"common/internal/domain/logger"
	"common/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

type commonRouter struct {
	R           *gin.Engine
	parent      *gin.RouterGroup
	tokenConfig config.TokenHashConfig
	logger      logger.Logger
}

type CommonRouter interface {
	SetAppValidationRoutes(handler *handler.AppValidationServiceHandler)
	SetAppTokenRoutes(handler *handler.AppTokenHandler)
	SetSkinRoutes(handler *handler.SkinHandler)
	SetConfigurationRoutes(handlers *handler.ConfigurationHandler)
	SetUserRoutes(handler *handler.UserHandler)
	SetInitAppValidtaionRoutes(handler *handler.AppValidationServiceHandler)
	SetDeviceRoutes(handler *handler.DeviceHandler)
	SetOrgRoutes(handler *handler.OrgHandler)

	GetEngine() *gin.Engine
}

func NewCommonRouter(serviceName string, tokenConfig config.TokenHashConfig, logger logger.Logger) CommonRouter {
	r := gin.Default()
	parent := r.Group("/" + serviceName)
	return &commonRouter{
		R:           r,
		parent:      parent,
		tokenConfig: tokenConfig,
		logger:      logger,
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
	client.Use(middleware.LoggingMiddleware(r.logger))
	//client.Use(middleware.AuthMiddleware(r.tokenConfig)) // JWT 적용
	client.POST("/refresh", handler.AppTokenRefresh)

}

func (r *commonRouter) SetSkinRoutes(handler *handler.SkinHandler) {
	client := r.parent.Group("/client/v1/skin-img")
	client.Use(middleware.LoggingMiddleware(r.logger))
	client.Use(middleware.AuthMiddleware(r.tokenConfig, r.logger)) // JWT 적용
	client.GET("", handler.GetSkinImage)

	server := r.parent.Group("/server/v1/skin-img")
	server.POST("", handler.PutSkinImg)
}

// 이 API도 service 처리.
// appValidation 했던것 처럼  config, skin을 아우르는  service생성하고 api의 전달되는 데이터에 따라서 조회하여 response하도록 수정.
func (r *commonRouter) SetConfigurationRoutes(handlers *handler.ConfigurationHandler) {
	// client := parent.Group("/client/v1/config-hash")
	// client.GET("/config-hash", handlers.ClientHandler.GetConfigHash)
}

func (r *commonRouter) SetUserRoutes(handler *handler.UserHandler) {
	// http://bookstack.ucware.local/books/neo-erd/page/5ef05에 등록된 사용자 등록 요청 (회원가입)
	client := r.parent.Group("/client/v1/user/register")
	client.Use(middleware.LoggingMiddleware(r.logger))
	client.POST("", handler.UserRegister)
	client.GET("/challenge", handler.GetUserRegisterChallenge)

}

func (r *commonRouter) SetInitAppValidtaionRoutes(handler *handler.AppValidationServiceHandler) {
	client := r.parent.Group("/client/v1/app-validation")
	client.Use(middleware.LoggingMiddleware(r.logger))
	client.GET("", handler.GetAppValidation)
}

func (r *commonRouter) SetDeviceRoutes(handler *handler.DeviceHandler) {
	server := r.parent.Group("/server/v1/device-init")
	server.POST("", handler.DeviceInit)
}

func (r *commonRouter) SetOrgRoutes(handler *handler.OrgHandler) {

}
