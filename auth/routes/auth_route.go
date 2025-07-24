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
			certification := v1.Group("/certification") // 도메인을 기준으로.
			{
				certification.POST("/login", authHandler.Login)
			}
		}

		sv1 := auth.Group("/sv1")
		{
			token := sv1.Group("/token") // 도메인을 기준으로
			{
				token.POST("/generate-app-token", serverHandler.GenerateAppToken)
				token.POST("/app-token-validation", serverHandler.AppTokenValidation)
				token.POST("/app-token-refresh", serverHandler.AppTokenRefresh)
			}
		}

	}

	return r
}
