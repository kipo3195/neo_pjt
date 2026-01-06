package di

import (
	"notificator/internal/application/orchestrator"
	"notificator/internal/application/usecase"
	"notificator/internal/delivery/handler"
	"notificator/internal/infrastructure/config"
)

func InitNotificatorServiceModule(chatRoom usecase.ChatRoomUsecase, socketSender usecase.SocketSenderUsecase, login usecase.LoginUsecase, websocketConfig config.WebsocketConnectionConfig) *handler.NotificatorServiceHandler {

	service := orchestrator.NewNotificatorService(chatRoom, socketSender, login)
	return handler.NewNotificatorServiceHandler(service, websocketConfig)

}
