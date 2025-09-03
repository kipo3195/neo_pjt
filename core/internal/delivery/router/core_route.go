package router

import (
	"core/internal/delivery/handler"

	"github.com/gin-gonic/gin"
)

func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	return r, r.Group("/" + serviceName)
}

func SetAppValidationRoute(parent *gin.RouterGroup, handler *handler.AppValidationHandler) {

	client := parent.Group("/client/v1/app-validation")
	client.POST("/validate", handler.ValidateApp)

}
