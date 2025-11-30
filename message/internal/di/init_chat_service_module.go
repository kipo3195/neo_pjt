package di

import (
	"message/internal/application/orchestrator"
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
)

func InitChatServiceModule(chat usecase.ChatUsecase, lineKey usecase.LineKeyUsecase, chatRoom usecase.ChatRoomUsecase) *handler.ChatServiceHandler {

	service := orchestrator.NewChatService(chat, lineKey, chatRoom)
	return handler.NewChatServiceHandler(service)

}
