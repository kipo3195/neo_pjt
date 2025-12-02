package orchestrator

import "notificator/internal/application/usecase"

type NotificatorService struct {
	Chat         usecase.ChatUsecase
	Note         usecase.NoteUsecase
	SocketSender usecase.SocketSenderUsecase
}

func NewNotificatorService(chat usecase.ChatUsecase, note usecase.NoteUsecase, socketSender usecase.SocketSenderUsecase) *NotificatorService {
	return &NotificatorService{
		Chat:         chat,
		Note:         note,
		SocketSender: socketSender,
	}
}
