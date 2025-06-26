package routes

import (
	"common/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(commonHandler *handlers.CommonHandler, serverHandler *handlers.ServerHandler, commonPubHandler *handlers.CommonPubHandler) *mux.Router {
	r := mux.NewRouter()

	// middleware가 필요없는 요청을 처리함.
	commonPub := r.PathPrefix("/common/pub").Subrouter()
	commonPub.HandleFunc("/app-validation", commonPubHandler.AppValidation).Methods("GET")

	commonV1 := r.PathPrefix("/common/v1").Subrouter()
	// TODO 클라이언트 middleware
	commonV1.Use(AuthMiddleware)
	commonV1.HandleFunc("/config-hash", commonHandler.GetConfigHash).Methods("GET")

	//----------------------------------------------------------------------------------------------------------------------------//

	// TODO 서버 middleware
	commonSV1 := r.PathPrefix("/common/sv1").Subrouter()
	commonSV1.HandleFunc("/device-init", serverHandler.DeviceInit).Methods("POST")

	return r
}
