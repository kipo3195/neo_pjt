package orchestrator

import "message/internal/application/usecase"

type ChatService struct {
	Chat    usecase.ChatUsecase
	LineKey usecase.LineKeyUsecase
}

func NewChatService(chat usecase.ChatUsecase, lineKey usecase.LineKeyUsecase) *ChatService {
	return &ChatService{
		Chat:    chat,
		LineKey: lineKey,
	}
}
