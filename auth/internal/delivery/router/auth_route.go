package router

import (
	"auth/internal/delivery/handler"
	"auth/internal/delivery/middleware"
	"auth/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

type authRouter struct {
	tokenConfig config.TokenHashConfig
	R           *gin.Engine
	parent      *gin.RouterGroup
}

type AuthRouter interface {
	SetCertificationRoutes(handler *handler.CertificationHandler)
	SetTokenRoutes(handler *handler.TokenHandler)
	SetUserAuthRoutes(handler *handler.UserAuthHandler)
	SetUserAuthServiceRoutes(handler *handler.UserAuthServiceHandler)
	SetUserAuthDeviceServiceRoutes(handler *handler.DeviceAuthServiceHandler)
	GetEngine() *gin.Engine
}

func (r *authRouter) GetEngine() *gin.Engine {
	return r.R
}

func NewAuthRouter(serviceName string, tokenConfig config.TokenHashConfig) AuthRouter {
	r := gin.Default()
	parent := r.Group("/" + serviceName)
	return &authRouter{
		tokenConfig: tokenConfig,
		parent:      parent,
		R:           r,
	}
}

// func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
// 	r := gin.Default()
// 	return r, r.Group("/" + serviceName)
// }

func (r *authRouter) SetCertificationRoutes(handler *handler.CertificationHandler) {
	client := r.parent.Group("/client/v1/certification")
	client.POST("/login", handler.Login)
}

func (r *authRouter) SetTokenRoutes(handler *handler.TokenHandler) {

	client := r.parent.Group("/client/v1/token")
	client.Use(middleware.AuthMiddleware(r.tokenConfig))
	client.POST("/re-issue-at", handler.AccessTokenReIssue)

	// 20250929 만약, 추후 at, rt refresh 로직이 들어간다면.. 메모리 로딩 - authTokenStorage도 refresh 필수.
	server := r.parent.Group("/server/v1/token")
	server.POST("/generate-app-token", handler.GenerateAppToken)
	server.POST("/app-token-validation", handler.AppTokenValidation)
	server.POST("/app-token-refresh", handler.AppTokenRefresh)

}

func (r *authRouter) SetUserAuthRoutes(handler *handler.UserAuthHandler) {
	client := r.parent.Group("/client/v1/user/auth")
	client.POST("/challenge", handler.GenerateAuthChallenge)
	//client.POST("/", handler.GetUserAuth) // 서비스로 이관

	server := r.parent.Group("/server/v1/user/auth/info/register")
	server.POST("/", handler.UserAuthInfoRegister)
}

func (r *authRouter) SetUserAuthServiceRoutes(handler *handler.UserAuthServiceHandler) {
	client := r.parent.Group("/client/v1/user/auth")
	client.POST("/", handler.UserAuthAndDeviceCheck)
}

func (r *authRouter) SetUserAuthDeviceServiceRoutes(handler *handler.DeviceAuthServiceHandler) {
	client := r.parent.Group("/client/v1/device")
	client.POST("/regist", handler.DeviceRegist)
	client.POST("/refresh", handler.DeviceRefresh)

}
