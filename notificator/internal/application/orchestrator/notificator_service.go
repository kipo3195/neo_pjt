package orchestrator

import "notificator/internal/application/usecase"

type NotificatorService struct {
	ChatRoom     usecase.ChatRoomUsecase
	SocketSender usecase.SocketSenderUsecase
	Login        usecase.LoginUsecase
}

func NewNotificatorService(chatRoom usecase.ChatRoomUsecase, socketSender usecase.SocketSenderUsecase, login usecase.LoginUsecase) *NotificatorService {
	return &NotificatorService{
		ChatRoom:     chatRoom,
		SocketSender: socketSender,
		Login:        login,
	}
}
