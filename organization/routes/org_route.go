package routes

import (
	"org/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(orgHandler *handlers.OrgHandler) *mux.Router {
	r := mux.NewRouter()
	// 메인 /org + 서브 라우터 활용

	/* 클라이언트가 호출하는 API */
	// v2 생성시 coreV2 := r.PathPrefix("/org/v2").Subrouter()..
	orgV1 := r.PathPrefix("/org/v1").Subrouter()

	orgV1.Use(AuthMiddleware)

	// 최초에 한하여 org의 모든 조직도 정보. 로컬 방식, response는 프로토콜 버퍼나 압축된 형태여야함
	orgV1.HandleFunc("/orgs", orgHandler.GetOrg).Methods("GET")

	// 요청하는 부서에 대한 조회. DB 방식, 최상위 포함.
	orgV1.HandleFunc("/departments", orgHandler.GetDept).Methods("GET")

	//----------------------------------------------------------------------------------------------------------------------------//

	/* 서버에서 호출하는 API */
	orgSV1 := r.PathPrefix("/org/sv1").Subrouter()

	// 서버 간 토큰 인증용 미들웨어
	// orgSV1.Use(ServerAuthMiddleware)

	// 부서 생성
	orgSV1.HandleFunc("/departments", orgHandler.ServerCreateDept).Methods("POST")

	// 삭제
	orgSV1.HandleFunc("/departments", orgHandler.ServerDeleteDept).Methods("DELETE")

	return r
}
