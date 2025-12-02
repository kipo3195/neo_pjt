package orchestrator

import "notificator/internal/application/usecase"

type NotificatorService struct {
	Chat         usecase.ChatUsecase
	Note         usecase.NoteUsecase
	SocketSender usecase.SocketSenderUsecase
	Login        usecase.LoginUsecase
}

func NewNotificatorService(chat usecase.ChatUsecase, note usecase.NoteUsecase, socketSender usecase.SocketSenderUsecase, login usecase.LoginUsecase) *NotificatorService {
	return &NotificatorService{
		Chat:         chat,
		Note:         note,
		SocketSender: socketSender,
		Login:        login,
	}
}
