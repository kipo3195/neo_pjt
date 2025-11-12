package di

import (
	"notificator/internal/application/orchestrator"
	"notificator/internal/application/usecase"
	"notificator/internal/delivery/handler"
)

func InitNotificatorServiceModule(chat usecase.ChatUsecase, note usecase.NoteUsecase) *handler.NotificatorServiceHandler {

	service := orchestrator.NewNotificatorService(chat, note)
	return handler.NewNotificatorServiceHandler(service)

}
