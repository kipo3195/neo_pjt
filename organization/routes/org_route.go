package routes

import (
	"org/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(orgHandler *handlers.OrgHandler) *mux.Router {
	r := mux.NewRouter()
	// 메인 /org + 서브 라우터 활용
	// v2 생성시 coreV2 := r.PathPrefix("/core/v2").Subrouter()..
	orgV1 := r.PathPrefix("/org/v1").Subrouter()

	orgV1.HandleFunc("/config", orgHandler.Test).Methods("POST")

	return r
}
