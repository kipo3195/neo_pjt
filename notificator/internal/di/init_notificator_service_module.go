package di

import (
	"notificator/internal/application/orchestrator"
	"notificator/internal/application/usecase"
	"notificator/internal/delivery/handler"
	"notificator/internal/infrastructure/config"
)

func InitNotificatorServiceModule(chat usecase.ChatUsecase, note usecase.NoteUsecase, socketSender usecase.SocketSenderUsecase, login usecase.LoginUsecase, websocketConfig config.WebsocketConnectionConfig) *handler.NotificatorServiceHandler {

	service := orchestrator.NewNotificatorService(chat, note, socketSender, login)
	return handler.NewNotificatorServiceHandler(service, websocketConfig)

}
