package router

import (
	appValidation "common/internal/domains/appValidation"

	"github.com/gin-gonic/gin"
)

func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	return r, r.Group("/" + serviceName)
}

func SetAppValidationRoutes(parent *gin.RouterGroup, handlers *appValidation.AppValidationHandlers) {

	server := parent.Group("/server/v1/app-validation")
	server.GET("/validate", handlers.ServerHandler.AppValidation)

}
