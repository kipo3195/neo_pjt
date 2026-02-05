package handler

import "file/internal/application/usecase"

type ChatFileHandler struct {
	usecase usecase.ChatFileUsecase
}

func NewChatFileHandler(usecase usecase.ChatFileUsecase) *ChatFileHandler {
	return &ChatFileHandler{
		usecase: usecase,
	}
}
