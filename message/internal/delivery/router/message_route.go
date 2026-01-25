package router

import (
	"message/internal/delivery/handler"
	"message/internal/delivery/middleware"
	"message/internal/domain/logger"
	"message/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

type messageRouter struct {
	tokenConfig config.TokenHashConfig
	R           *gin.Engine
	parent      *gin.RouterGroup
	logger      logger.Logger
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
	SetChatLineServiceRoutes(handler *handler.ChatLineServiceHandler)
}

func (r *messageRouter) GetEngine() *gin.Engine {
	return r.R
}

func NewMessageRouter(serviceName string, tokenConfig config.TokenHashConfig, logger logger.Logger) MessageRouter {
	r := gin.Default()

	// 해당 서비스의 모든 API 요청에 대한 로깅 적용
	// parent 밑에서 로깅 미들웨어 적용시 /wrong-path로 접속했을때 그룹 매칭에 실패하여 미들웨어가 아예 타지 않기 때문.
	r.Use(middleware.LoggingMiddleware(logger))
	parent := r.Group("/" + serviceName)
	return &messageRouter{
		parent:      parent,
		R:           r,
		tokenConfig: tokenConfig,
		logger:      logger,
	}
}

func (r *messageRouter) SetLineKeyRoutes(handler *handler.LineKeyHandler) {

	client := r.parent.Group("/client/v1/line-key")
	client.Use(middleware.AuthMiddleware(r.tokenConfig, r.logger))
	client.GET("", handler.GetLineKey)
}

func (r *messageRouter) SetChatRoutes(handler *handler.ChatHandler) {

}

func (r *messageRouter) SetChatServiceRoutes(handler *handler.ChatServiceHandler) {
	client := r.parent.Group("/client/v1/chat")
	client.Use(middleware.AuthMiddleware(r.tokenConfig, r.logger))
	client.POST("", handler.SendChat)
	client.POST("/read", handler.ReadChat)
}

func (r *messageRouter) SetNoteRoutes(handler *handler.NoteHandler) {

	client := r.parent.Group("/client/v1/note")
	client.Use(middleware.AuthMiddleware(r.tokenConfig, r.logger))
	client.POST("", handler.SendNote)

}

func (r *messageRouter) SetOtpRoutes(handler *handler.OtpHandler) {
	server := r.parent.Group("/server/v1/otp")
	server.POST("/regist", handler.OtpKeyRegist)

	client := r.parent.Group("/client/v1/otp")
	client.Use(middleware.AuthMiddleware(r.tokenConfig, r.logger))
	client.POST("/info", handler.GetMyOtpInfo)
}

func (r *messageRouter) SetChatRoomRoutes(handler *handler.ChatRoomHandler) {

}

func (r *messageRouter) SetChatRoomTitleRoutes(handler *handler.ChatRoomTitleHandler) {

	client := r.parent.Group("client/v1/chat/room/title")
	client.Use(middleware.AuthMiddleware(r.tokenConfig, r.logger))
	client.POST("", handler.UpdateChatRoomTitle)
	client.DELETE("", handler.DeleteChatRoomTitle)

}

func (r *messageRouter) SetChatRoomServiceRoutes(handler *handler.ChatRoomServiceHandler) {
	client := r.parent.Group("/client/v1/chat/room")
	client.Use(middleware.AuthMiddleware(r.tokenConfig, r.logger))
	client.POST("", handler.CreateChatRoom)
	client.POST("/detail", handler.GetChatRoomDetail)
	client.POST("/list", handler.GetChatRoomList)
	client.POST("/update-date", handler.GetChatRoomUpdateDate)
	client.POST("/member/read-date", handler.GetChatRoomMemberReadDate)
	client.POST("/my", handler.GetChatRoomMy)
	client.POST("/sync, ")
}

func (r *messageRouter) SetChatLineServiceRoutes(handler *handler.ChatLineServiceHandler) {

	client := r.parent.Group("/client/v1/chat/line")
	client.Use(middleware.AuthMiddleware(r.tokenConfig, r.logger))
	client.POST("/event", handler.GetChatLineEvent)
}
