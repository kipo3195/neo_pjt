package routes

import (
	"core/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(handlers *handlers.CoreHandlers) *mux.Router {
	r := mux.NewRouter()
	// 메인 /core + 서브 라우터 활용
	// v2 생성시 coreV2 := r.PathPrefix("/core/v2").Subrouter()..
	coreV1 := r.PathPrefix("/core/v1").Subrouter()

	// TODO client 미들웨어
	coreV1.HandleFunc("/app-validation", handlers.Core.GetAppValidation).Methods("POST")

	//----------------------------------------------------------------------------------------------------------------------------//

	return r
}
