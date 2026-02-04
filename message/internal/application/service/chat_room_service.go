package service

import "message/internal/application/usecase"

type ChatRoomService struct {
	ChatRoom       usecase.ChatRoomUsecase
	LineKey        usecase.LineKeyUsecase
	Chat           usecase.ChatUsecase
	ChatRoomFixed  usecase.ChatRoomFixedUsecase
	ChatRoomTitle  usecase.ChatRoomTitleUsecase
	ChatRoomConfig usecase.ChatRoomConfigUsecase
}

func NewChatRoomService(chatRoom usecase.ChatRoomUsecase, lineKey usecase.LineKeyUsecase, chat usecase.ChatUsecase, chatRoomFixed usecase.ChatRoomFixedUsecase,
	chatRoomTitle usecase.ChatRoomTitleUsecase, chatRoomConfig usecase.ChatRoomConfigUsecase) *ChatRoomService {

	return &ChatRoomService{
		ChatRoom:       chatRoom,
		LineKey:        lineKey,
		Chat:           chat,
		ChatRoomFixed:  chatRoomFixed,
		ChatRoomTitle:  chatRoomTitle,
		ChatRoomConfig: chatRoomConfig,
	}
}
