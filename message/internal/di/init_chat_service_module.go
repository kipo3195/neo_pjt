package di

import (
	"message/internal/application/orchestrator"
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
)

func InitChatServiceModule(chat usecase.ChatUsecase, lineKey usecase.LineKeyUsecase) *handler.ChatServiceHandler {

	service := orchestrator.NewChatService(chat, lineKey)
	return handler.NewChatServiceHandler(service)

}
