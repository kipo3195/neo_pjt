package routes

import (
	"auth/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(authHandler *handlers.AuthHandler, serverHandler *handlers.ServerHandler) *gin.Engine {

	r := gin.Default()

	auth := r.Group("/auth")
	{

		v1 := auth.Group("/v1")
		{
			v1.POST("/login", authHandler.Login)
		}

		sv1 := auth.Group("/sv1")
		{
			sv1.POST("/generate-app-token", serverHandler.GenerateAppToken)
			sv1.POST("/app-token-validation", serverHandler.AppTokenValidation)
		}

	}

	return r
}
