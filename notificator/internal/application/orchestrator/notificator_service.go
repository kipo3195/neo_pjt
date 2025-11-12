package orchestrator

import "notificator/internal/application/usecase"

type NotificatorService struct {
	Chat usecase.ChatUsecase
	Note usecase.NoteUsecase
}

func NewNotificatorService(chat usecase.ChatUsecase, note usecase.NoteUsecase) *NotificatorService {
	return &NotificatorService{
		Chat: chat,
		Note: note,
	}
}
