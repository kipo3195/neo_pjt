package routes

import (
	"auth/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(authHandler *handlers.AuthHandler, serverHandler *handlers.ServerHandler) *mux.Router {
	r := mux.NewRouter()

	authV1 := r.PathPrefix("/auth/v1").Subrouter()

	// 최초 인증시 토큰 발급
	authV1.HandleFunc("/generate-device-token", authHandler.GenerateDeviceToken).Methods("POST")

	// 사용자 인증
	authV1.HandleFunc("/login", authHandler.Login).Methods("POST")

	//----------------------------------------------------------------------------------------------------------------------------//

	// 장시간 인증 X 상태에서 fore 왔을때 앱 토큰 검증 common서버로 부터 호출 받음. 마지막 발급 키(limit 1)를 기준으로 체크
	authSV1 := r.PathPrefix("/auth/sv1").Subrouter()
	authSV1.HandleFunc("/app-token-validation", serverHandler.AppTokenValidation).Methods("POST")

	return r
}
