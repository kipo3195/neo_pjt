package router

import (
	"test/internal/adapter/http/handler"

	"github.com/gin-gonic/gin"
)

type testRouter struct {
	R      *gin.Engine
	parent *gin.RouterGroup
}

type TestRouter interface {
	GetEngine() *gin.Engine
	SetLoginRoutes(handler *handler.LoginHandler)
	SetChatRoutes(handler *handler.ChatHandler)
}

func (r *testRouter) GetEngine() *gin.Engine {
	return r.R
}

func NewTestRouter(serviceName string) TestRouter {
	r := gin.Default()
	parent := r.Group("/" + serviceName)

	return &testRouter{
		parent: parent,
		R:      r,
	}
}

func (r *testRouter) SetLoginRoutes(handler *handler.LoginHandler) {
	client := r.parent.Group("/client/v1/login")
	client.POST("", handler.PutLogin)
}

func (r *testRouter) SetChatRoutes(handler *handler.ChatHandler) {
	client := r.parent.Group("/client/v1/chat")
	client.POST("", handler.PutChatEvent)
}
