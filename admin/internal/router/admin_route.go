package router

import (
	"admin/internal/domains/orgDeptUsers"
	"admin/internal/domains/orgDepts"
	"admin/internal/domains/orgFile"
	"admin/internal/domains/skinImg"

	"github.com/gin-gonic/gin"
)

func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	return r, r.Group("/" + serviceName)
}

func SetOrgDeptUsersRoutes(parent *gin.RouterGroup, handlers *orgDeptUsers.OrgDeptUsersHandlers) {

	client := parent.Group("/client/v1/org/dept/users")
	client.GET("/", handlers.ClientHandler.GetUsers) // 전체 조회 (특정 조회도 필요하다면? )
	client.POST("/", handlers.ClientHandler.RegisterUsers)
	client.PUT("/", handlers.ClientHandler.UpdateUsers)
	client.DELETE("/", handlers.ClientHandler.DeleteUsers)

}

func SetOrgDeptsRoutes(parent *gin.RouterGroup, handlers *orgDepts.OrgDeptsHandlers) {

	client := parent.Group("/client/v1/org/depts")
	client.GET("/", handlers.ClientHandler.GetDepts) // 전체 조회 (특정 조회도 필요하다면? )
	client.POST("/", handlers.ClientHandler.RegisterDepts)
	client.PUT("/", handlers.ClientHandler.UpdateDepts)
	client.DELETE("/", handlers.ClientHandler.DeleteDepts)

}

func SetOrgFilesRoutes(parent *gin.RouterGroup, handlers *orgFile.OrgFileHandlers) {

	client := parent.Group("/client/v1/org/file")
	client.POST("/", handlers.ClientHandler.CreateOrgFile)
	client.GET("/", handlers.ClientHandler.GetOrgFile)

}

func SetSkinRoutes(parent *gin.RouterGroup, handlers *skinImg.SkinImgHandlers) {

	client := parent.Group("/client/v1/skinImg")
	client.POST("/", handlers.ClientHandler.CreateSkinImg)

}
