package routes

import (
	"admin/handlers"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(handlers *handlers.AdminHandlers) *gin.Engine {
	r := gin.Default()

	r.Use(TimeoutMiddleware(5 * time.Second))

	admin := r.Group("/admin")
	{
		v1 := admin.Group("/v1")
		{
			v1.POST("/org/departments", handlers.Org.CreateDept)
			v1.GET("/org/departments", handlers.Org.GetDept)
			v1.PUT("/org/departments", handlers.Org.UpdateDept)
			v1.DELETE("/org/departments", handlers.Org.DeleteDept)

			v1.POST("/org/users", handlers.Org.CreateUser)
			v1.GET("/org/users", handlers.Org.GetUser)
			v1.PUT("/org/users", handlers.Org.UpdateUser)
			v1.DELETE("/org/users", handlers.Org.DeleteUser)

			v1.POST("/org/file", handlers.Org.CreateOrgFile)
			v1.GET("/org/file", handlers.Org.GetOrgFile)

			v1.POST("/common/skin-img", handlers.Common.CreateSkinImg)

		}
	}

	return r
}
