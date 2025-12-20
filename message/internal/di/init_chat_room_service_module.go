package di

import (
	"message/internal/application/orchestrator"
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
)

func InitChatRoomServiceModule(chatRoom usecase.ChatRoomUsecase, lineKey usecase.LineKeyUsecase, chat usecase.ChatUsecase, chatRoomFixed usecase.ChatRoomFixedUsecase,
	chatRoomTitle usecase.ChatRoomTitleUsecase, chatRoomConfig usecase.ChatRoomConfigUsecase) *handler.ChatRoomServiceHandler {

	service := orchestrator.NewChatRoomService(chatRoom, lineKey, chat, chatRoomFixed, chatRoomTitle, chatRoomConfig)
	return handler.NewChatRoomServiceHandler(service)

}
