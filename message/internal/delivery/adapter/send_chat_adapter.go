package adapter

import (
	"message/internal/application/usecase/input"
	"message/internal/delivery/dto/chatService"
)

func MakeSendChatInput(sendUserHash string, eventType string, lineKey string, chatRoom chatService.ChatRoomData, chatLine chatService.ChatLineData) input.SendChatInput {

	room := input.ChatRoomInput{
		RoomKey:  chatRoom.RoomKey,
		RoomType: chatRoom.RoomType,
	}

	line := input.ChatLineInput{
		SendUserHash: sendUserHash,
		EventType:    eventType,
		LineKey:      lineKey,
		Contents:     chatLine.Contents,
	}

	return input.SendChatInput{
		ChatRoom: room,
		ChatLine: line,
	}
}
