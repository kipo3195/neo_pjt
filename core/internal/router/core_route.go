package router

import (
	"core/internal/domains/appValidation"

	"github.com/gin-gonic/gin"
)

func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	return r, r.Group("/" + serviceName)
}

func SetAppValidationRoutes(parent *gin.RouterGroup, handlers *appValidation.AppValidationHandlers) {

	client := parent.Group("/client/v1/app-validation")
	client.POST("/validate", handlers.ClientHandler.ValidateApp)

}
