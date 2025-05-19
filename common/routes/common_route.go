package routes

import (
	"common/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(commonHandler *handlers.CommonHandler) *mux.Router {
	r := mux.NewRouter()

	commonV1 := r.PathPrefix("/common/v1").Subrouter()

	commonV1.HandleFunc("/device-init", commonHandler.DeviceInit).Methods("POST")
	commonV1.HandleFunc("/get-config", commonHandler.GetConfig).Methods("POST")

	return r
}
