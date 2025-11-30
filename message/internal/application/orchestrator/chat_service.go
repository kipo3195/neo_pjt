package orchestrator

import "message/internal/application/usecase"

type ChatService struct {
	Chat     usecase.ChatUsecase
	LineKey  usecase.LineKeyUsecase
	ChatRoom usecase.ChatRoomUsecase
}

func NewChatService(chat usecase.ChatUsecase, lineKey usecase.LineKeyUsecase, chatRoom usecase.ChatRoomUsecase) *ChatService {
	return &ChatService{
		Chat:     chat,
		LineKey:  lineKey,
		ChatRoom: chatRoom,
	}
}
