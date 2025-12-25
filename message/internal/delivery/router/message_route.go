package router

import (
	"message/internal/delivery/handler"
	"message/internal/delivery/middleware"
	"message/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

type messageRouter struct {
	tokenConfig config.TokenHashConfig
	R           *gin.Engine
	parent      *gin.RouterGroup
}

type MessageRouter interface {
	GetEngine() *gin.Engine
	SetLineKeyRoutes(handler *handler.LineKeyHandler)
	SetChatRoutes(handler *handler.ChatHandler)
	SetChatServiceRoutes(handler *handler.ChatServiceHandler)
	SetNoteRoutes(handler *handler.NoteHandler)
	SetOtpRoutes(handler *handler.OtpHandler)
	SetChatRoomRoutes(handler *handler.ChatRoomHandler)
	SetChatRoomServiceRoutes(handler *handler.ChatRoomServiceHandler)
	SetChatRoomTitleRoutes(handler *handler.ChatRoomTitleHandler)
}

func (r *messageRouter) GetEngine() *gin.Engine {
	return r.R
}

func NewMessageRouter(serviceName string, tokenConfig config.TokenHashConfig) MessageRouter {
	r := gin.Default()
	parent := r.Group("/" + serviceName)
	return &messageRouter{
		parent:      parent,
		R:           r,
		tokenConfig: tokenConfig,
	}
}

func (r *messageRouter) SetLineKeyRoutes(handler *handler.LineKeyHandler) {

	client := r.parent.Group("/client/v1/line-key")
	client.Use(middleware.AuthMiddleware(r.tokenConfig))
	client.GET("/", handler.GetLineKey)
}

func (r *messageRouter) SetChatRoutes(handler *handler.ChatHandler) {

}

func (r *messageRouter) SetChatServiceRoutes(handler *handler.ChatServiceHandler) {

	client := r.parent.Group("/client/v1/chat")
	client.Use(middleware.AuthMiddleware(r.tokenConfig))
	client.POST("", handler.SendChat)

}

func (r *messageRouter) SetNoteRoutes(handler *handler.NoteHandler) {

	client := r.parent.Group("/client/v1/note")
	client.Use(middleware.AuthMiddleware(r.tokenConfig))
	client.POST("", handler.SendNote)

}

func (r *messageRouter) SetOtpRoutes(handler *handler.OtpHandler) {
	server := r.parent.Group("/server/v1/otp")
	server.POST("/regist", handler.OtpKeyRegist)

	client := r.parent.Group("/client/v1/otp")
	client.Use(middleware.AuthMiddleware(r.tokenConfig))
	client.POST("/info", handler.GetMyOtpInfo)
}

func (r *messageRouter) SetChatRoomRoutes(handler *handler.ChatRoomHandler) {

}

func (r *messageRouter) SetChatRoomTitleRoutes(handler *handler.ChatRoomTitleHandler) {

	client := r.parent.Group("client/v1/chat/room/title")
	client.Use(middleware.AuthMiddleware(r.tokenConfig))
	client.POST("", handler.UpdateChatRoomTitle)
	client.DELETE("", handler.DeleteChatRoomTitle)

}

func (r *messageRouter) SetChatRoomServiceRoutes(handler *handler.ChatRoomServiceHandler) {
	client := r.parent.Group("/client/v1/chat/room")
	client.Use(middleware.AuthMiddleware(r.tokenConfig))
	client.POST("", handler.CreateChatRoom)
	client.POST("/detail", handler.GetChatRoomDetail)
	client.POST("/list", handler.GetChatRoomList)
}
