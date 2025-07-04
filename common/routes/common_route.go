package routes

import (
	"common/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(handlers *handlers.CommonHandlers) *mux.Router {
	r := mux.NewRouter()

	// 토큰을 검증하지 않는 로직
	commonPub := r.PathPrefix("/common/pub").Subrouter()

	commonPub.HandleFunc("/app-validation", handlers.Public.AppValidation).Methods("GET")

	commonV1 := r.PathPrefix("/common/v1").Subrouter()

	// TODO 클라이언트 middleware
	commonV1.Use(AuthMiddleware)

	// configHash 검증
	commonV1.HandleFunc("/config-hash", handlers.Common.GetConfigHash).Methods("GET")

	// 스킨 이미지 요청
	commonV1.HandleFunc("/skin-img", handlers.Common.GetSkinImage).Methods("GET")

	//----------------------------------------------------------------------------------------------------------------------------//

	// TODO 서버 middleware
	commonSV1 := r.PathPrefix("/common/sv1").Subrouter()

	// core
	commonSV1.HandleFunc("/device-init", handlers.Server.DeviceInit).Methods("POST")

	// admin 스킨 파일 업로드
	commonSV1.HandleFunc("/skin-img", handlers.Server.PutSkinImg).Methods("POST")

	return r
}
