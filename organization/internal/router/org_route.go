package router

import (
	"org/internal/domains/department"
	"org/internal/domains/org"
	"org/internal/domains/user"

	"github.com/gin-gonic/gin"
)

func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	return r, r.Group("/" + serviceName)
}

func SetDepartmentRoutes(parent *gin.RouterGroup, handlers *department.DepartmentHandlers) {
	client := parent.Group("/client/v1/department")
	client.GET("/", handlers.ClientHandler.GetDept)

	server := parent.Group("/server/v1/department")
	server.POST("/", handlers.ServerHandler.CreateDept)
	server.DELETE("/", handlers.ServerHandler.DeleteDept)
	server.POST("/user", handlers.ServerHandler.DeleteDept)
	server.DELETE("/user", handlers.ServerHandler.DeleteDept)

}

func SetOrgRoute(parent *gin.RouterGroup, handlers *org.OrgHandlers) {

	client := parent.Group("/client/v1/org")
	client.GET("/hash", handlers.ClientHandler.GetOrgHash)
	client.GET("/data", handlers.ClientHandler.GetOrgData)

	server := parent.Group("/server/v1/org")
	server.POST("/file", handlers.ServerHandler.CreateOrgFile)

}

func SetUserRoute(parent *gin.RouterGroup, handlers *user.UserHandlers) {
	client := parent.Group("/client/v1/user")
	client.GET("/my-info", handlers.ClientHandler.GetMyInfo)
	client.GET("/info", handlers.ClientHandler.GetUserInfo)

}
