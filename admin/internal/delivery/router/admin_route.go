package router

import (
	"admin/internal/delivery/handler"

	"github.com/gin-gonic/gin"
)

func SetDefaultRoutes(serviceName string) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	return r, r.Group("/" + serviceName)
}

func SetOrgDeptUserRoutes(parent *gin.RouterGroup, handler *handler.OrgDeptUserHandler) {

	client := parent.Group("/client/v1/org/dept/user")
	client.GET("/", handler.GetDeptUser) // 전체 조회 (특정 조회도 필요하다면? )
	client.POST("/", handler.RegistDeptUser)
	client.PUT("/", handler.UpdateDeptUser)
	client.DELETE("/", handler.DeleteDeptUser)

}

func SetOrgDeptRoutes(parent *gin.RouterGroup, handler *handler.OrgDeptHandler) {

	client := parent.Group("/client/v1/org/dept")
	client.GET("/", handler.GetDept) // 전체 조회 (특정 조회도 필요하다면? )
	client.POST("/", handler.RegisterDept)
	client.PUT("/", handler.UpdateDept)
	client.DELETE("/", handler.DeleteDept)

}

func SetOrgFileRoutes(parent *gin.RouterGroup, handler *handler.OrgFileHandler) {

	client := parent.Group("/client/v1/org/file")
	client.POST("/", handler.CreateOrgFile)
	client.GET("/", handler.GetOrgFile)

}

func SetSkinRoutes(parent *gin.RouterGroup, handler *handler.SkinImgHandler) {

	client := parent.Group("/client/v1/skinImg")
	client.POST("/", handler.CreateSkinImg)

}
