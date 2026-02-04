package mapper

import (
	"message/internal/adapter/http/dto/chatService"
	"message/internal/application/usecase/input"
)

func MakeSendChatInput(sendUserHash string, lineKey string, sendDate string, eventType string, chatSession string, chatRoom chatService.ChatRoomData, chatLine chatService.ChatLineData, transactionId string) input.SendChatInput {

	room := input.ChatRoomInput{
		RoomKey:    chatRoom.RoomKey,
		RoomType:   chatRoom.RoomType,
		SecretFlag: chatRoom.SecretFlag,
	}

	line := input.ChatLineInput{
		SendUserHash:  sendUserHash,
		LineKey:       lineKey,
		TargetLineKey: chatLine.TargetLineKey,
		Contents:      chatLine.Contents,
		Cmd:           chatLine.Cmd,
		SendDate:      sendDate,
	}

	return input.SendChatInput{
		ChatRoom:      room,
		ChatLine:      line,
		ChatSession:   chatSession,
		EventType:     eventType,
		TransactionId: transactionId,
	}
}
