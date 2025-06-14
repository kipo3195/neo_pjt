package routes

import (
	"admin/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(ah *handlers.AdminHandlers) *mux.Router {
	r := mux.NewRouter()

	adminV1 := r.PathPrefix("/admin/v1").Subrouter()

	// 인증 미들웨어, 타임아웃 미들웨어 적용
	//adminV1.Use(AuthMiddleware)
	adminV1.Use(TimeoutMiddleware)

	// /admin/v1/org/departments
	// 부서 CRUD
	adminV1.HandleFunc("/org/departments", ah.AdminOrgHandler.CreateDept).Methods("POST")
	adminV1.HandleFunc("/org/departments", ah.AdminOrgHandler.GetDept).Methods("GET") // 부서조회 - 부서 + 사용자
	adminV1.HandleFunc("/org/departments", ah.AdminOrgHandler.UpdateDept).Methods("PUT")
	adminV1.HandleFunc("/org/departments", ah.AdminOrgHandler.DeleteDept).Methods("DELETE")

	// /admin/v1/org/users
	// 사용자 CUD -> 관리자에서 조회하는건 DB 기반으로 처리해도 되지않을까?
	adminV1.HandleFunc("/org/users", ah.AdminOrgHandler.CreateUser).Methods("POST")
	adminV1.HandleFunc("/org/users", ah.AdminOrgHandler.GetUser).Methods("GET") // 사용자 - 등급, 다국어 등..
	adminV1.HandleFunc("/org/users", ah.AdminOrgHandler.UpdateUser).Methods("PUT")
	adminV1.HandleFunc("/org/users", ah.AdminOrgHandler.DeleteUser).Methods("DELETE")

	// 현재 DB를 기준으로 org를 파일로 만드는 API
	adminV1.HandleFunc("/org/file", ah.AdminOrgHandler.CreateOrgFile).Methods("POST") // 생성
	adminV1.HandleFunc("/org/file", ah.AdminOrgHandler.GetOrgFile).Methods("GET")     // 조회

	return r
}
