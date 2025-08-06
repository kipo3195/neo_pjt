package routes

import (
	"common/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(handlers *handlers.CommonHandlers) *gin.Engine {

	r := gin.Default()

	common := r.Group("/common")
	{

		// 일부 API에 middleware 적용
		v1 := common.Group("/v1")
		v1.Use(AuthMiddleware)

		{
			v1.GET("/config-hash", handlers.Common.GetConfigHash)

		}

		sv1 := common.Group("/sv1")
		{
			//sv1.POST("/device-init", handlers.Server.DeviceInit)
		}
	}

	// // 토큰을 검증하지 않는 로직
	// commonPub := r.PathPrefix("/common/pub").Subrouter()

	// commonPub.HandleFunc("/app-validation", handlers.Public.AppValidation).Methods("GET")

	// commonPub.HandleFunc("/app-token-refresh", handlers.Public.AppTokenRefresh).Methods("POST")

	// commonV1 := r.PathPrefix("/common/v1").Subrouter()

	// // TODO 클라이언트 middleware
	// commonV1.Use(AuthMiddleware)

	// // configHash 검증
	// commonV1.HandleFunc("/config-hash", handlers.Common.GetConfigHash).Methods("GET")

	// // 스킨 이미지 요청
	// commonV1.HandleFunc("/skin-img", handlers.Common.GetSkinImage).Methods("GET")

	// //----------------------------------------------------------------------------------------------------------------------------//

	// // TODO 서버 middleware
	// commonSV1 := r.PathPrefix("/common/sv1").Subrouter()

	// // core
	// commonSV1.HandleFunc("/device-init", handlers.Server.DeviceInit).Methods("POST")

	// // admin 스킨 파일 업로드
	// commonSV1.HandleFunc("/skin-img", handlers.Server.PutSkinImg).Methods("POST")

	return r
}
