package routes

import (
	"admin/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(adminHandler *handlers.AdminHandler) *mux.Router {
	r := mux.NewRouter()

	adminV1 := r.PathPrefix("/admin/v1").Subrouter()

	// 인증 미들웨어, 타임아웃 미들웨어 적용
	adminV1.Use(AuthMiddleware)
	adminV1.Use(TimeoutMiddleware)

	// 부서 CRUD
	adminV1.HandleFunc("/departments", adminHandler.CreateDept).Methods("POST")
	adminV1.HandleFunc("/departments", adminHandler.GetDept).Methods("GET")
	adminV1.HandleFunc("/departments", adminHandler.UpdateDept).Methods("PUT")
	adminV1.HandleFunc("/departments", adminHandler.DeleteDept).Methods("DELETE")

	return r
}
