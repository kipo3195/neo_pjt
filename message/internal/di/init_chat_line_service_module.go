package di

import (
	"message/internal/application/orchestrator"
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
)

func InitChatLineServiceModule(chat usecase.ChatUsecase, chatRoom usecase.ChatRoomUsecase) *handler.ChatLineServiceHandler {

	service := orchestrator.NewChatLineService(chat, chatRoom)
	return handler.NewChatLineServiceHandler(service)

}
