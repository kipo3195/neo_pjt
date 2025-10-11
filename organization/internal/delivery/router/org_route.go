package router

import (
	"org/internal/delivery/handler"
	"org/internal/delivery/middleware"
	"org/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	return r, r.Group("/" + serviceName)
}

func SetDepartmentRoutes(parent *gin.RouterGroup, handler *handler.DepartmentHandler, tokenConfig config.TokenHashConfig) {
	clientApi := parent.Group("/client/v1/department")
	clientApi.Use(middleware.AuthMiddleware(tokenConfig))
	clientApi.GET("/", handler.GetDept) //

	serverApi := parent.Group("/server/v1/department")
	serverApi.Use(middleware.ServerAuthMiddleware())
	serverApi.POST("/", handler.CreateDept)
	serverApi.DELETE("/", handler.DeleteDept)
	serverApi.POST("/user", handler.CreateDeptUser)
	serverApi.DELETE("/user", handler.DeleteDeptUser)

}

func SetOrgRoute(parent *gin.RouterGroup, handler *handler.OrgHandler, tokenConfig config.TokenHashConfig) {

	clientApi := parent.Group("/client/v1/org")
	clientApi.Use(middleware.AuthMiddleware(tokenConfig))
	clientApi.GET("/hash", handler.GetOrgHash)
	clientApi.GET("/data", handler.GetOrgData)

	serverApi := parent.Group("/server/v1/org")
	//serverApi.Use(middleware.ServerAuthMiddleware())
	serverApi.POST("/file", handler.CreateOrgFile)

}

func SetUserRoute(parent *gin.RouterGroup, handler *handler.UserHandler, tokenConfig config.TokenHashConfig) {
	clientApi := parent.Group("/client/v1/user")

	clientApi.Use(middleware.AuthMiddleware(tokenConfig))
	clientApi.GET("/my-info", handler.GetMyInfo)
	clientApi.GET("/info", handler.GetUserInfo)
}

// 더미데이터 생성 Service /////////////////////////////
func SetDummyDataServiceRoute(parent *gin.RouterGroup, handler *handler.DummyDataServiceHandler) {

	user := parent.Group("/test/v1/user")
	user.POST("/init/service-user/", handler.InitServiceUser)
	user.POST("/init/user-detail/", handler.InitUserDetail)
	user.POST("/init/user-multi-lang", handler.InitUserMultiLang)

	department := parent.Group("/test/v1/department")
	department.POST("/init/works-dept", handler.InitWorksDept)
	department.POST("/init/works-dept-multi-lang", handler.InitWorksDeptMultiLang)
	department.POST("/init/works-dept-user", handler.InitWorksDeptUser)

}
