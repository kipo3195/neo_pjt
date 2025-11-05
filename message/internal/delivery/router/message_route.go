package router

import (
	"github.com/gin-gonic/gin"
)

type messageRouter struct {
	R      *gin.Engine
	parent *gin.RouterGroup
}

type MessageRouter interface {
	GetEngine() *gin.Engine
}

func (r *messageRouter) GetEngine() *gin.Engine {
	return r.R
}

func NewMessageRouter(serviceName string) MessageRouter {
	r := gin.Default()
	parent := r.Group("/" + serviceName)
	return &messageRouter{
		parent: parent,
		R:      r,
	}
}
