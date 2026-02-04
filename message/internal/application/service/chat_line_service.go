package orchestrator

import "message/internal/application/usecase"

type ChatLineService struct {
	Chat     usecase.ChatUsecase
	ChatRoom usecase.ChatRoomUsecase
	ChatFile usecase.ChatFileUsecase
}

func NewChatLineService(chat usecase.ChatUsecase, chatRoom usecase.ChatRoomUsecase, chatFile usecase.ChatFileUsecase) *ChatLineService {
	return &ChatLineService{
		Chat:     chat,
		ChatRoom: chatRoom,
		ChatFile: chatFile,
	}
}
