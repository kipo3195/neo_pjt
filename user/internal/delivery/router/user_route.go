package router

import (
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
