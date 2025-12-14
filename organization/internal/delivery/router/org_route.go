package router

import (
	"org/internal/delivery/handler"
	"org/internal/delivery/middleware"
	"org/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

type orgRouter struct {
	R      *gin.Engine
	parent *gin.RouterGroup
}

type OrgRouter interface {
	SetDepartmentRoutes(handler *handler.DepartmentHandler, tokenConfig config.TokenHashConfig)
	SetOrgRoute(handler *handler.OrgHandler, tokenConfig config.TokenHashConfig)
	SetUserRoute(handler *handler.UserHandler, tokenConfig config.TokenHashConfig)
	SetDummyDataServiceRoute(handler *handler.DummyDataServiceHandler)
	GetEngine() *gin.Engine
}

func NewOrgRoute(serviceName string) OrgRouter {

	r := gin.Default()
	parent := r.Group("/" + serviceName)

	return &orgRouter{
		R:      r,
		parent: parent,
	}
}

func (r *orgRouter) GetEngine() *gin.Engine {
	return r.R
}

func (r *orgRouter) SetDepartmentRoutes(handler *handler.DepartmentHandler, tokenConfig config.TokenHashConfig) {
	clientApi := r.parent.Group("/client/v1/department")
	clientApi.Use(middleware.AuthMiddleware(tokenConfig))
	clientApi.GET("/", handler.GetDept) //

	serverApi := r.parent.Group("/server/v1/department")
	serverApi.Use(middleware.ServerAuthMiddleware())
	serverApi.POST("/", handler.CreateDept)
	serverApi.DELETE("/", handler.DeleteDept)
	serverApi.POST("/user", handler.CreateDeptUser)
	serverApi.DELETE("/user", handler.DeleteDeptUser)

}

func (r *orgRouter) SetOrgRoute(handler *handler.OrgHandler, tokenConfig config.TokenHashConfig) {

	clientApi := r.parent.Group("/client/v1/org")
	clientApi.Use(middleware.AuthMiddleware(tokenConfig))
	clientApi.GET("/hash", handler.GetOrgHash)
	clientApi.GET("/data", handler.GetOrgData)

	serverApi := r.parent.Group("/server/v1/org")
	//serverApi.Use(middleware.ServerAuthMiddleware())
	serverApi.POST("/batch", handler.RegistOrgBatch) // batch 서비스에서 호출한 연동된 조직정보 데이터
	serverApi.POST("/file", handler.CreateOrgFile)   // admin 서비스에서 호출한 현재 기준으로 json 생성 요청

}

func (r *orgRouter) SetUserRoute(handler *handler.UserHandler, tokenConfig config.TokenHashConfig) {
	clientApi := r.parent.Group("/client/v1/user")

	clientApi.Use(middleware.AuthMiddleware(tokenConfig))
	clientApi.GET("/my-info", handler.GetMyInfo)
	clientApi.POST("/info", handler.GetUserInfo)
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
