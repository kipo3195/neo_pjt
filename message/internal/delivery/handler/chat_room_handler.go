package handler

import "message/internal/application/usecase"

type ChatRoomHandler struct {
	usecase usecase.ChatRoomUsecase
}

func NewChatRoomHandler(usecase usecase.ChatRoomUsecase) *ChatRoomHandler {
	return &ChatRoomHandler{
		usecase: usecase,
	}
}
