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
	client.POST("/img", handler.UploadProfileImg)
	client.GET("/img", handler.GetProfileImg)
	client.DELETE("/img", handler.DeleteProfileImg) // 기본 이미지로 변경

	client.POST("/msg", handler.RegistProfileMsg)

}
