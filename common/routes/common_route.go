package routes

import (
	"common/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(commonHandler *handlers.CommonHandler) *mux.Router {
	r := mux.NewRouter()

	commonV1 := r.PathPrefix("/common/v1").Subrouter()

	// TODO 클라이언트 middleware
	commonV1.HandleFunc("/get-config", commonHandler.GetConfig).Methods("POST")

	//----------------------------------------------------------------------------------------------------------------------------//

	// TODO 서버 middleware
	commonV1.HandleFunc("/device-init", commonHandler.DeviceInit).Methods("POST")

	return r
}
