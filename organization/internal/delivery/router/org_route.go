package router

import (
	"org/internal/delivery/handler"
	"org/internal/delivery/middleware"
	"org/internal/domain/logger"
	"org/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

type orgRouter struct {
	R           *gin.Engine
	parent      *gin.RouterGroup
	tokenConfig config.TokenHashConfig
	logger      logger.Logger
}

type OrgRouter interface {
	SetDepartmentRoutes(handler *handler.DepartmentHandler)
	SetOrgRoute(handler *handler.OrgHandler)
	SetUserRoute(handler *handler.UserHandler)
	SetDummyDataServiceRoute(handler *handler.DummyDataServiceHandler)
	SetOrgBatchServiceRoute(handler *handler.OrgBatchServiceHandler)
	SetOrgUserServiceRoute(handler *handler.OrgUserServiceHandler)
	GetEngine() *gin.Engine
}

func NewOrgRoute(serviceName string, tokenConfig config.TokenHashConfig, logger logger.Logger) OrgRouter {

	r := gin.Default()

	// 해당 서비스의 모든 API 요청에 대한 로깅 적용
	// parent 밑에서 로깅 미들웨어 적용시 /wrong-path로 접속했을때 그룹 매칭에 실패하여 미들웨어가 아예 타지 않기 때문.
	r.Use(middleware.LoggingMiddleware(logger))
	parent := r.Group("/" + serviceName)

	return &orgRouter{
		R:           r,
		parent:      parent,
		tokenConfig: tokenConfig,
		logger:      logger,
	}
}

func (r *orgRouter) GetEngine() *gin.Engine {
	return r.R
}

func (r *orgRouter) SetDepartmentRoutes(handler *handler.DepartmentHandler) {
	client := r.parent.Group("/client/v1/department")
	client.Use(middleware.AuthMiddleware(r.tokenConfig, r.logger))
	client.GET("/", handler.GetDept) //

	serverApi := r.parent.Group("/server/v1/department")
	serverApi.Use(middleware.ServerAuthMiddleware())
	serverApi.POST("/", handler.CreateDept)
	serverApi.DELETE("/", handler.DeleteDept)
	serverApi.POST("/user", handler.CreateDeptUser)
	serverApi.DELETE("/user", handler.DeleteDeptUser)

}

func (r *orgRouter) SetOrgRoute(handler *handler.OrgHandler) {

	client := r.parent.Group("/client/v1/org")
	client.Use(middleware.AuthMiddleware(r.tokenConfig, r.logger))
	client.GET("/hash", handler.GetOrgHash)
	client.GET("/data", handler.GetOrgData)

	serverApi := r.parent.Group("/server/v1/org")
	//serverApi.Use(middleware.ServerAuthMiddleware())
	serverApi.POST("/file", handler.CreateOrgFile) // admin 서비스에서 호출한 현재 기준으로 json 생성 요청

}

func (r *orgRouter) SetUserRoute(handler *handler.UserHandler) {
	// clientApi := r.parent.Group("/client/v1/user")

	// clientApi.Use(middleware.AuthMiddleware(tokenConfig))
	//clientApi.GET("/my-info", handler.GetMyInfo)
	//clientApi.POST("/info", handler.GetUserInfo)
}

// 더미데이터 생성 Service /////////////////////////////
func (r *orgRouter) SetDummyDataServiceRoute(handler *handler.DummyDataServiceHandler) {

	user := r.parent.Group("/test/v1/user")
	user.POST("/init/service-user/", handler.InitServiceUser)
	user.POST("/init/user-detail/", handler.InitUserDetail)
	user.POST("/init/user-multi-lang", handler.InitUserMultiLang)

	department := r.parent.Group("/test/v1/department")
	department.POST("/init/works-dept", handler.InitWorksDept)
	department.POST("/init/works-dept-multi-lang", handler.InitWorksDeptMultiLang)
	department.POST("/init/works-dept-user", handler.InitWorksDeptUser)

}

func (r *orgRouter) SetOrgBatchServiceRoute(handler *handler.OrgBatchServiceHandler) {

	server := r.parent.Group("/server/v1/org/batch")
	server.POST("/", handler.RegistOrgBatchData)

}

func (r *orgRouter) SetOrgUserServiceRoute(handler *handler.OrgUserServiceHandler) {

	client := r.parent.Group("/client/v1/user")
	client.Use(middleware.AuthMiddleware(r.tokenConfig, r.logger))
	client.GET("/my-info", handler.GetMyInfo)
	client.POST("/info", handler.GetUserInfo)

}
