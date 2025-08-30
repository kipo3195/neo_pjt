package router

import (
	"org/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	return r, r.Group("/" + serviceName)
}

func SetDepartmentRoutes(parent *gin.RouterGroup, handler *handler.DepartmentHandler) {
	clientApi := parent.Group("/client/v1/department")
	clientApi.GET("/", handler.GetDept) //

	serverApi := parent.Group("/server/v1/department")
	serverApi.POST("/", handler.CreateDept)
	serverApi.DELETE("/", handler.DeleteDept)
	serverApi.POST("/user", handler.CreateDeptUser)
	serverApi.DELETE("/user", handler.DeleteDeptUser)

}

func SetOrgRoute(parent *gin.RouterGroup, handler *handler.OrgHandler) {

	client := parent.Group("/client/v1/org")
	client.GET("/hash", handler.GetOrgHash)
	client.GET("/data", handler.GetOrgData)

	server := parent.Group("/server/v1/org")
	server.POST("/file", handler.CreateOrgFile)

}

func SetUserRoute(parent *gin.RouterGroup, handler *handler.UserHandler) {
	client := parent.Group("/client/v1/user")
	client.GET("/my-info", handler.GetMyInfo)
	client.GET("/info", handler.GetUserInfo)

}
