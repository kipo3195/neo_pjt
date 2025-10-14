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
	client := parent.Group("/client/v1/token")
	client.POST("/re-issue-at", handler.AccessTokenRefresh)

	// 20250929 만약, 추후 at, rt refresh 로직이 들어간다면.. 메모리 로딩 - authTokenStorage도 refresh 필수.
	server := parent.Group("/server/v1/token")
	server.POST("/generate-app-token", handler.GenerateAppToken)
	server.POST("/app-token-validation", handler.AppTokenValidation)
	server.POST("/app-token-refresh", handler.AppTokenRefresh)

}

func SetUserAuthRoutes(parent *gin.RouterGroup, handler *handler.UserAuthHandler) {
	client := parent.Group("/client/v1/user/auth")
	client.POST("/challenge", handler.GenerateAuthChallenge)
	//client.POST("/", handler.GetUserAuth) // 서비스로 이관

	server := parent.Group("/server/v1/user/auth/info/register")
	server.POST("/", handler.UserAuthInfoRegister)
}

func SetUserAuthServiceRoutes(parent *gin.RouterGroup, handler *handler.UserAuthServiceHandler) {
	client := parent.Group("/client/v1/user/auth")
	client.POST("/", handler.UserAuthAndDeviceCheck)
}

func SetDeviceRoutes(parent *gin.RouterGroup, handler *handler.DeviceHandler) {
	client := parent.Group("/client/v1/device")
	client.POST("/regist", handler.DeviceRegist)
}
