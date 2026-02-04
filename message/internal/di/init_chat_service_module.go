package di

import (
	"message/internal/adapter/http/handler"
	"message/internal/application/service"
	"message/internal/application/usecase"
)

func InitChatServiceModule(chat usecase.ChatUsecase, lineKey usecase.LineKeyUsecase, chatRoom usecase.ChatRoomUsecase) *handler.ChatServiceHandler {

	service := service.NewChatService(chat, lineKey, chatRoom)
	return handler.NewChatServiceHandler(service)

}
