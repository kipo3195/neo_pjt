package router

import (
	"user/internal/delivery/handler"
	"user/internal/delivery/middleware"
	"user/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

type userRouter struct {
	tokenConfig config.TokenHashConfig
	R           *gin.Engine
	parent      *gin.RouterGroup
}

type UserRouter interface {
	GetEngine() *gin.Engine
	SetProfileRoutes(handler *handler.ProfileHandler)
	SetUserDetailRoutes(handler *handler.UserDetailHandler)
	SetUserInfoServiceRoutes(handler *handler.UserInfoServiceHandler)
}

func (r *userRouter) GetEngine() *gin.Engine {
	return r.R
}

func NewUserRouter(serviceName string, tokenConfig config.TokenHashConfig) UserRouter {
	r := gin.Default()
	parent := r.Group("/" + serviceName)
	return &userRouter{
		tokenConfig: tokenConfig,
		parent:      parent,
		R:           r,
	}
}

func (r *userRouter) SetProfileRoutes(handler *handler.ProfileHandler) {
	client := r.parent.Group("/client/v1/profile")
	client.Use(middleware.AuthMiddleware(r.tokenConfig))

	client.POST("/img/upload", handler.UploadProfileImg)
	client.POST("/img", handler.GetProfileImg)
	client.DELETE("/img", handler.DeleteProfileImg) // 기본 이미지로 변경

	client.POST("/msg", handler.RegistProfileMsg)
}

func (r *userRouter) SetUserDetailRoutes(handler *handler.UserDetailHandler) {

	// 사용자의 ID가 아닌 HASH 정보로 요청해야하므로 부담스러운 GET보다는 POST로 요청
	client := r.parent.Group("/client/v1/detail")
	client.Use(middleware.AuthMiddleware(r.tokenConfig))
	//client.GET("/my-info", handler.GetMyDetailInfo) // 정보 조회
	//client.POST("/info", handler.GetUserDetailInfo) // 정보 조회

	// 생각해 봐야 할것은 endponit의 형식. detail을 한번에 수정하는지, 부분적으로 수정하는지
	// uri의 데이터를 분기 -> /client/v1/detail/name, /client/v1/detail/email...
	// 전체 일괄 분기 /client/v1/detail에 POST 방식

}

func (r *userRouter) SetUserInfoServiceRoutes(handler *handler.UserInfoServiceHandler) {
	client := r.parent.Group("/client/v1/info")
	client.Use(middleware.AuthMiddleware(r.tokenConfig))
	client.GET("/my", handler.GetMyDetailInfo)      // 내 정보 조회
	client.POST("/user", handler.GetUserDetailInfo) // 정보 조회
}
