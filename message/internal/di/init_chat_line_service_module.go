package di

import (
	"message/internal/adapter/http/handler"
	"message/internal/adapter/rpc/grpcHandler"
	"message/internal/application/service"
	"message/internal/application/usecase"
)

type ChatLineServiceModule struct {
	Handler     *handler.ChatLineServiceHandler
	GrpcHandler *grpcHandler.ChatLineServiceHandler
}

func InitChatLineServiceModule(chat usecase.ChatUsecase, chatRoom usecase.ChatRoomUsecase, chatFile usecase.ChatFileUsecase) ChatLineServiceModule {

	service := service.NewChatLineService(chat, chatRoom, chatFile)
	return ChatLineServiceModule{
		handler.NewChatLineServiceHandler(service),
		grpcHandler.NewChatLineServiceHandler(service),
	}
}
