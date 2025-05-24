package routes

import (
	"auth/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(authHandler *handlers.AuthHandler) *mux.Router {
	r := mux.NewRouter()

	authV1 := r.PathPrefix("/auth/v1").Subrouter()
	authV1.HandleFunc("/auth", authHandler.GetAuth).Methods("POST")
	authV1.HandleFunc("/generate-device-token", authHandler.GenerateDeviceToken).Methods("POST")

	return r
}
