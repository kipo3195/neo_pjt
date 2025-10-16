package router

import (
	"core/internal/delivery/handler"

	"github.com/gin-gonic/gin"
)

type coreRouter struct {
	R      *gin.Engine
	parent *gin.RouterGroup
}

type CoreRouter interface {
	SetAppValidationRoute(handler *handler.AppValidationHandler)
	GetEngine() *gin.Engine
}

func NewCoreRouter(serviceName string) CoreRouter {
	r := gin.Default()
	parent := r.Group("/" + serviceName)
	return &coreRouter{
		R:      r,
		parent: parent,
	}
}

func (r *coreRouter) GetEngine() *gin.Engine {
	return r.R
}

func (r *coreRouter) SetAppValidationRoute(handler *handler.AppValidationHandler) {

	client := r.parent.Group("/client/v1/app-validation")
	client.POST("/validate", handler.ValidateApp)

}
