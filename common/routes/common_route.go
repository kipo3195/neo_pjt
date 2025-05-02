package routes

import (
	"common/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(commonHandler *handlers.CommonHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/getConfig", commonHandler.GetConfig).Methods("POST")
	return r
}
