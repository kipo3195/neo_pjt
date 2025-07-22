package routes

import (
	"core/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(handlers *handlers.CoreHandlers) *gin.Engine {

	r := gin.Default()

	core := r.Group("/core")
	{
		// 	/core/v1
		v1 := core.Group("/v1")
		{
			v1.POST("/app-validation", handlers.Core.GetAppValidation)
		}

	}

	return r
}
