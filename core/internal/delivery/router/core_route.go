package router

import (
	"core/internal/delivery/handler"
	"core/internal/domain/logger"

	"core/internal/delivery/middleware"

	"github.com/gin-gonic/gin"
)

type coreRouter struct {
	R      *gin.Engine
	parent *gin.RouterGroup
	logger logger.Logger
}

type CoreRouter interface {
	SetAppValidationRoute(handler *handler.AppValidationHandler)
	GetEngine() *gin.Engine
}

func NewCoreRouter(serviceName string, logger logger.Logger) CoreRouter {
	r := gin.Default()
	parent := r.Group("/" + serviceName)
	return &coreRouter{
		R:      r,
		parent: parent,
		logger: logger,
	}
}

func (r *coreRouter) GetEngine() *gin.Engine {
	return r.R
}

func (r *coreRouter) SetAppValidationRoute(handler *handler.AppValidationHandler) {

	client := r.parent.Group("/client/v1/app-validation")
	client.Use(middleware.LoggingMiddleware(r.logger))
	client.POST("/validate", handler.ValidateApp)

}
