package di

import (
	"message/internal/adapter/http/handler"
	"message/internal/application/service"
	"message/internal/application/usecase"
)

func InitChatLineServiceModule(chat usecase.ChatUsecase, chatRoom usecase.ChatRoomUsecase, chatFile usecase.ChatFileUsecase) *handler.ChatLineServiceHandler {

	service := service.NewChatLineService(chat, chatRoom, chatFile)
	return handler.NewChatLineServiceHandler(service)

}
