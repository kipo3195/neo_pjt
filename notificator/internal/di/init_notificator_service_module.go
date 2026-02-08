package di

import (
	"notificator/internal/adapter/http/handler"
	"notificator/internal/application/service"
	"notificator/internal/application/usecase"
	"notificator/internal/infrastructure/config"
)

func InitNotificatorServiceModule(chatRoom usecase.ChatRoomUsecase, socketSender usecase.SocketSenderUsecase, login usecase.LoginUsecase, websocketConfig config.WebsocketConnectionConfig) *handler.NotificatorServiceHandler {

	service := service.NewNotificatorService(chatRoom, socketSender, login)
	return handler.NewNotificatorServiceHandler(service, websocketConfig)

}
