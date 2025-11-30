package adapter

import (
	"message/internal/application/usecase/input"
	"message/internal/delivery/dto/chatService"
)

func MakeSendChatInput(sendUserHash string, lineKey string, sendDate string, eventType string, chatSession string, chatRoom chatService.ChatRoomData, chatLine chatService.ChatLineData) input.SendChatInput {

	room := input.ChatRoomInput{
		RoomKey:  chatRoom.RoomKey,
		RoomType: chatRoom.RoomType,
	}

	line := input.ChatLineInput{
		SendUserHash: sendUserHash,
		LineKey:      lineKey,
		Contents:     chatLine.Contents,
		Cmd:          chatLine.Cmd,
		SendDate:     sendDate,
	}

	return input.SendChatInput{
		ChatRoom:    room,
		ChatLine:    line,
		ChatSession: chatSession,
		EventType:   eventType,
	}
}
