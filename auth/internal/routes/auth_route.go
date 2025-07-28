package routes

import (
	certificationClientHandler "auth/internal/domains/certification/handlers/client"

	"github.com/gin-gonic/gin"
)

func SetDefaultRoutes(serviceName string) *gin.RouterGroup {
	r := gin.Default()
	return r.Group("/" + serviceName)
}

func SetLoginRoute(parent *gin.RouterGroup, handler *certificationClientHandler.CertificationHandler) {

	group := parent.Group("/client/v1/certification")
	group.POST("/login", handler.Login)

}

func SetupServerTokenRoute(parent *gin.RouterGroup) *gin.Engine {

	group := parent.Group("/server/v1/token")
	group.POST("/generate-app-token", serverTokenHandler.GenerateAppToken)
	group.POST("/app-token-validation", serverTokenHandler.AppTokenValidation)
	group.POST("/app-token-refresh", serverTokenHandler.AppTokenRefresh)

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
