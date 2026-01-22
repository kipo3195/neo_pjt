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

	// 해당 서비스의 모든 API 요청에 대한 로깅 적용
	// parent 밑에서 로깅 미들웨어 적용시 /wrong-path로 접속했을때 그룹 매칭에 실패하여 미들웨어가 아예 타지 않기 때문.
	r.Use(middleware.LoggingMiddleware(logger))
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
	client.POST("/validate", handler.ValidateApp)

}
