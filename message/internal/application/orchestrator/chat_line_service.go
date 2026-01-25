package orchestrator

import "message/internal/application/usecase"

type ChatLineService struct {
	Chat     usecase.ChatUsecase
	ChatRoom usecase.ChatRoomUsecase
}

func NewChatLineService(chat usecase.ChatUsecase, chatRoom usecase.ChatRoomUsecase) *ChatLineService {
	return &ChatLineService{
		Chat:     chat,
		ChatRoom: chatRoom,
	}
}
