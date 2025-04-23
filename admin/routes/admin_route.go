package routes

import (
	"admin/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(ruleHandler *handlers.RuleHandler, authHandler *handlers.AuthHandler) *mux.Router {
	r := mux.NewRouter()
	return r
}
