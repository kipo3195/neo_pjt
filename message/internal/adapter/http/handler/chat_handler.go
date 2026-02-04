package handler

import "message/internal/application/usecase"

type ChatHandler struct {
	usecase usecase.ChatUsecase
}

func NewChatHandler(uc usecase.ChatUsecase) *ChatHandler {
	return &ChatHandler{
		usecase: uc,
	}
}
