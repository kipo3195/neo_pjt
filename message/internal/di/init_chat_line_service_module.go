package di

import (
	"message/internal/application/orchestrator"
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
)

func InitChatLineServiceModule(chat usecase.ChatUsecase, chatRoom usecase.ChatRoomUsecase, chatFile usecase.ChatFileUsecase) *handler.ChatLineServiceHandler {

	service := orchestrator.NewChatLineService(chat, chatRoom, chatFile)
	return handler.NewChatLineServiceHandler(service)

}
