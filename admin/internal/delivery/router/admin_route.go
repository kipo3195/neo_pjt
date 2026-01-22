package router

import (
	"admin/internal/delivery/handler"
	"admin/internal/delivery/middleware"
	"admin/internal/domain/logger"

	"github.com/gin-gonic/gin"
)

type adminRouter struct {
	R      *gin.Engine
	parent *gin.RouterGroup
	logger logger.Logger
}

type AdminRouter interface {
	SetOrgDeptUserRoutes(handler *handler.OrgDeptUserHandler)
	SetOrgDeptRoutes(handler *handler.OrgDeptHandler)
	SetOrgFileRoutes(handler *handler.OrgFileHandler)
	SetSkinRoutes(handler *handler.SkinImgHandler)
	SetServiceUserRoutes(handler *handler.ServiceUserHandler)
	SetServiceUserAuthRegisterServiceRoutes(handler *handler.ServiceUserAuthRegisterHandler)
	GetEngine() *gin.Engine
}

func NewAdminRouter(serviceName string, logger logger.Logger) AdminRouter {
	r := gin.Default()
	parent := r.Group("/" + serviceName)
	return &adminRouter{
		R:      r,
		parent: parent,
		logger: logger, // 인증 미들웨어에 주입 하기 위함
	}
}

func (r *adminRouter) GetEngine() *gin.Engine {
	return r.R
}

func (r *adminRouter) SetOrgDeptUserRoutes(handler *handler.OrgDeptUserHandler) {

	client := r.parent.Group("/client/v1/org/dept/user")
	client.Use(middleware.LoggingMiddleware(r.logger))
	client.GET("", handler.GetDeptUser) // 전체 조회 (특정 조회도 필요하다면? )
	client.POST("", handler.RegistDeptUser)
	client.PUT("", handler.UpdateDeptUser)
	client.DELETE("", handler.DeleteDeptUser)

}

func (r *adminRouter) SetOrgDeptRoutes(handler *handler.OrgDeptHandler) {

	client := r.parent.Group("/client/v1/org/dept")
	client.Use(middleware.LoggingMiddleware(r.logger))
	client.GET("", handler.GetDept) // 전체 조회 (특정 조회도 필요하다면? )
	client.POST("", handler.RegisterDept)
	client.PUT("", handler.UpdateDept)
	client.DELETE("", handler.DeleteDept)

}

func (r *adminRouter) SetOrgFileRoutes(handler *handler.OrgFileHandler) {

	client := r.parent.Group("/client/v1/org/file")
	client.Use(middleware.LoggingMiddleware(r.logger))
	client.POST("", handler.CreateOrgFile)
	client.GET("", handler.GetOrgFile)

}

func (r *adminRouter) SetSkinRoutes(handler *handler.SkinImgHandler) {

	client := r.parent.Group("/client/v1/skinImg")
	client.Use(middleware.LoggingMiddleware(r.logger))
	client.POST("", handler.CreateSkinImg)

}

func (r *adminRouter) SetServiceUserRoutes(handler *handler.ServiceUserHandler) {
	// 20251217 인증정보 처리로직 추가로 인해서 SetServiceUserAuthRegisterServiceRoutes로 이관
	// client := r.parent.Group("/client/v1/serviceUser")
	//client.POST("/", handler.RegistServiceUser)
}

func (r *adminRouter) SetServiceUserAuthRegisterServiceRoutes(handler *handler.ServiceUserAuthRegisterHandler) {
	client := r.parent.Group("/client/v1/serviceUser")
	client.Use(middleware.LoggingMiddleware(r.logger))
	client.POST("", handler.RegistServiceUser)
}
