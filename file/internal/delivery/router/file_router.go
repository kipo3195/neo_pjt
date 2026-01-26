package router

import (
	"file/internal/delivery/handler"
	"file/internal/delivery/middleware"
	"file/internal/domain/logger"
	"file/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

type fileRouter struct {
	tokenConfig config.TokenHashConfig
	R           *gin.Engine
	parent      *gin.RouterGroup
	logger      logger.Logger
}

type FileRouter interface {
	GetEngine() *gin.Engine
	SetFileUrlRoutes(handler *handler.FileUrlHandler)
}

func (r *fileRouter) GetEngine() *gin.Engine {
	return r.R
}

func NewFileRouter(serviceName string, tokenConfig config.TokenHashConfig, logger logger.Logger) FileRouter {
	r := gin.Default()

	// 해당 서비스의 모든 API 요청에 대한 로깅 적용
	// parent 밑에서 로깅 미들웨어 적용시 /wrong-path로 접속했을때 그룹 매칭에 실패하여 미들웨어가 아예 타지 않기 때문.
	r.Use(middleware.LoggingMiddleware(logger))
	parent := r.Group("/" + serviceName)
	return &fileRouter{
		tokenConfig: tokenConfig,
		parent:      parent,
		R:           r,
		logger:      logger,
	}
}

func (r *fileRouter) SetFileUrlRoutes(handler *handler.FileUrlHandler) {

	client := r.parent.Group("/client/v1/file-url")
	client.Use(middleware.AuthMiddleware(r.tokenConfig, r.logger))

	client.POST("", handler.CreateFileUrl)
	client.POST("/upload-end", handler.FileUrlUploadEnd)
}
