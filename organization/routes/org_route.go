package routes

import (
	"org/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(orgHandler *handlers.OrgHandler) *mux.Router {
	r := mux.NewRouter()
	// 메인 /org + 서브 라우터 활용
	// v2 생성시 coreV2 := r.PathPrefix("/org/v2").Subrouter()..
	orgV1 := r.PathPrefix("/org/v1").Subrouter()

	// 최초에 한하여 org의 모든 조직도 정보. 로컬 방식, response는 프로토콜 버퍼나 압축된 형태여야함
	orgV1.HandleFunc("/orgs", orgHandler.GetOrg).Methods("GET")

	// 요청하는 부서에 대한 조회. DB 방식, 최상위 포함.
	orgV1.HandleFunc("/departments", orgHandler.GetDept).Methods("GET")

	return r
}
