package di

import (
	"message/internal/adapter/http/handler"
	"message/internal/application/service"
	"message/internal/application/usecase"
)

func InitChatRoomServiceModule(chatRoom usecase.ChatRoomUsecase, lineKey usecase.LineKeyUsecase, chat usecase.ChatUsecase, chatRoomFixed usecase.ChatRoomFixedUsecase,
	chatRoomTitle usecase.ChatRoomTitleUsecase, chatRoomConfig usecase.ChatRoomConfigUsecase) *handler.ChatRoomServiceHandler {

	service := service.NewChatRoomService(chatRoom, lineKey, chat, chatRoomFixed, chatRoomTitle, chatRoomConfig)
	return handler.NewChatRoomServiceHandler(service)

}
