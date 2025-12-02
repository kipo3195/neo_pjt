package di

import (
	"notificator/internal/application/orchestrator"
	"notificator/internal/application/usecase"
	"notificator/internal/delivery/handler"
)

func InitNotificatorServiceModule(chat usecase.ChatUsecase, note usecase.NoteUsecase, socketSender usecase.SocketSenderUsecase) *handler.NotificatorServiceHandler {

	service := orchestrator.NewNotificatorService(chat, note, socketSender)
	return handler.NewNotificatorServiceHandler(service)

}
