package handlers

import "message/usecases"

type ChatHandler struct {
	usecase usecases.ChatUsecase
}

func NewChatHandler(uc usecases.ChatUsecase) *ChatHandler {
	return &ChatHandler{usecase: uc}
}
