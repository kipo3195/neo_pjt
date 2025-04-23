package routes

import (
	"auth/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(authHandler *handlers.AuthHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/auth", authHandler.GetAuth).Methods("POST")
	return r
}
