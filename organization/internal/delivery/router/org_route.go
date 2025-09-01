package router

import (
	"org/internal/delivery/handler"
	"org/internal/delivery/middleware"

	"github.com/gin-gonic/gin"
)

func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	return r, r.Group("/" + serviceName)
}

func SetDepartmentRoutes(parent *gin.RouterGroup, handler *handler.DepartmentHandler) {
	clientApi := parent.Group("/client/v1/department")
	clientApi.Use(middleware.AuthMiddleware())
	clientApi.GET("/", handler.GetDept) //

	serverApi := parent.Group("/server/v1/department")
	serverApi.Use(middleware.ServerAuthMiddleware())
	serverApi.POST("/", handler.CreateDept)
	serverApi.DELETE("/", handler.DeleteDept)
	serverApi.POST("/user", handler.CreateDeptUser)
	serverApi.DELETE("/user", handler.DeleteDeptUser)

}

func SetOrgRoute(parent *gin.RouterGroup, handler *handler.OrgHandler) {

	clientApi := parent.Group("/client/v1/org")
	clientApi.Use(middleware.AuthMiddleware())
	clientApi.GET("/hash", handler.GetOrgHash)
	clientApi.GET("/data", handler.GetOrgData)

	serverApi := parent.Group("/server/v1/org")
	serverApi.Use(middleware.ServerAuthMiddleware())
	serverApi.POST("/file", handler.CreateOrgFile)

}

func SetUserRoute(parent *gin.RouterGroup, handler *handler.UserHandler) {
	clientApi := parent.Group("/client/v1/user")
	clientApi.Use(middleware.AuthMiddleware())
	clientApi.GET("/my-info", handler.GetMyInfo)
	clientApi.GET("/info", handler.GetUserInfo)
}
