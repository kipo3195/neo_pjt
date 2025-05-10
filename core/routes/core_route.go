package routes

import (
	"core/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(coreHandler *handlers.CoreHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/v{version:[0-9]+}/getValidation", coreHandler.GetValidation).Methods("POST")
	return r
}
