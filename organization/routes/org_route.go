package routes

import (
	"org/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(orgHandler *handlers.OrgHandler, userHandler *handlers.UserHandler, serverHandler *handlers.ServerHandler) *mux.Router {
	// 메인 /org + 서브 라우터 활용
	r := mux.NewRouter()
	/* 클라이언트가 호출하는 API */
	// v2 생성시 coreV2 := r.PathPrefix("/org/v2").Subrouter()..
	orgV1 := r.PathPrefix("/org/v1").Subrouter()

	// 인증 처리 미들웨어
	orgV1.Use(AuthMiddleware)

	// org_hash의 배열 형태로 서버에 요청, 현재 서버의 hash와 비교해서 file, event response
	orgV1.HandleFunc("/orgs/hash", orgHandler.GetOrgHash).Methods("GET")

	// 요청하는 타입(이벤트, 파일)에 따라 조직도 response
	orgV1.HandleFunc("/orgs/data", orgHandler.GetOrgData).Methods("GET")

	// 요청하는 부서에 대한 조회. DB 방식, 최상위 포함.
	orgV1.HandleFunc("/departments", orgHandler.GetDept).Methods("GET")

	// 내 정보 조회
	orgV1.HandleFunc("/user/my-info", userHandler.GetMyInfo).Methods("GET")

	// 요청하는 사용자 정보 조회
	orgV1.HandleFunc("/user/info", userHandler.GetUserInfo).Methods("GET")

	//----------------------------------------------------------------------------------------------------------------------------//

	/* 서버에서 호출하는 API */
	orgSV1 := r.PathPrefix("/org/sv1").Subrouter()

	// 서버 간 토큰 인증용 미들웨어
	// orgSV1.Use(ServerAuthMiddleware)

	// 부서 생성
	orgSV1.HandleFunc("/departments", serverHandler.ServerCreateDept).Methods("POST")
	// 삭제
	orgSV1.HandleFunc("/departments", serverHandler.ServerDeleteDept).Methods("DELETE")

	// 부서에 사용자 추가
	orgSV1.HandleFunc("/departments/user", serverHandler.ServerCreateDeptUser).Methods("POST")
	// 삭제
	orgSV1.HandleFunc("/departments/user", serverHandler.ServerDeleteDeptUser).Methods("DELETE")

	// 현재 기준으로 org 파일 및 DB 저장
	orgSV1.HandleFunc("/org/file", serverHandler.ServerCreateOrgFile).Methods("POST")

	return r
}
