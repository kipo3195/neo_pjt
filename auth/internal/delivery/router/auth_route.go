package router

import (
	"auth/internal/delivery/handler"
	"auth/internal/delivery/middleware"
	"auth/internal/domain/logger"
	"auth/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

type authRouter struct {
	tokenConfig config.TokenHashConfig
	R           *gin.Engine
	parent      *gin.RouterGroup
	logger      logger.Logger
}

type AuthRouter interface {
	SetTokenRoutes(handler *handler.TokenHandler)
	SetUserAuthRoutes(handler *handler.UserAuthHandler)
	SetUserAuthServiceRoutes(handler *handler.UserAuthServiceHandler)
	SetUserAuthDeviceServiceRoutes(handler *handler.DeviceAuthServiceHandler)
	SetDeviceRoutes(handler *handler.DeviceHandler)
	GetEngine() *gin.Engine
}

func (r *authRouter) GetEngine() *gin.Engine {
	return r.R
}

func NewAuthRouter(serviceName string, tokenConfig config.TokenHashConfig, logger logger.Logger) AuthRouter {
	r := gin.Default()
	parent := r.Group("/" + serviceName)
	return &authRouter{
		tokenConfig: tokenConfig,
		parent:      parent,
		R:           r,
		logger:      logger,
	}
}

// func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
// 	r := gin.Default()
// 	return r, r.Group("/" + serviceName)
// }

func (r *authRouter) SetTokenRoutes(handler *handler.TokenHandler) {

	client := r.parent.Group("/client/v1/token")
	client.Use(middleware.LoggingMiddleware(r.logger))
	// AuthWithoutExpMiddleware at의 만료시간은 체크하지않고 사용자 아이디만 파싱처리하기 위함
	// 20251018 주석 처리한 이유 : AT를 클라이언트에서 읽어버린 경우 대비
	//client.Use(middleware.AuthWithoutExpMiddleware(r.tokenConfig))
	client.POST("/re-issue-at", handler.AccessTokenReIssue)

	// 20250929 만약, 추후 at, rt refresh 로직이 들어간다면.. 메모리 로딩 - authTokenStorage도 refresh 필수.
	server := r.parent.Group("/server/v1/token")
	server.POST("/generate-app-token", handler.GenerateAppToken)
	server.POST("/app-token-validation", handler.AppTokenValidation)
	server.POST("/app-token-refresh", handler.AppTokenRefresh)

}

func (r *authRouter) SetUserAuthRoutes(handler *handler.UserAuthHandler) {
	client := r.parent.Group("/client/v1/user/auth")
	client.Use(middleware.LoggingMiddleware(r.logger))
	client.POST("/challenge", handler.GenerateAuthChallenge)
	//client.POST("/", handler.GetUserAuth) // 서비스로 이관

	// common, admin을 통한 사용자 인증 정보등록 API
	// 20251217 배열로 변경
	server := r.parent.Group("/server/v1/user/auth/info/register")
	server.POST("/", handler.UserAuthInfoRegister)
}

func (r *authRouter) SetUserAuthServiceRoutes(handler *handler.UserAuthServiceHandler) {
	client := r.parent.Group("/client/v1/user/auth")
	client.Use(middleware.LoggingMiddleware(r.logger))
	client.POST("/", handler.UserAuthAndDeviceCheck)
}

func (r *authRouter) SetUserAuthDeviceServiceRoutes(handler *handler.DeviceAuthServiceHandler) {
	client := r.parent.Group("/client/v1/device")
	client.Use(middleware.LoggingMiddleware(r.logger))
	client.POST("/regist", handler.DeviceRegist)
	client.POST("/refresh", handler.DeviceRefresh)
}

func (r *authRouter) SetDeviceRoutes(handler *handler.DeviceHandler) {
	client := r.parent.Group("/client/v1/device")
	client.Use(middleware.LoggingMiddleware(r.logger))
	client.Use(middleware.AuthMiddleware(r.tokenConfig, r.logger))
	client.GET("/my-info", handler.GetMyDeviceInfo)
}
