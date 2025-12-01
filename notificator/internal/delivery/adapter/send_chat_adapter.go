package adapter

import (
	"notificator/internal/application/usecase/input"
	"notificator/internal/application/usecase/output"
)

func MakeSendChatInput(eventType string, chatSession string, roomData output.ChatRoomDataOutput, lineData output.ChatLineDataOutput) input.SendChatInput {

	chatRoomDataInput := input.ChatRoomDataInput{
		RoomKey:    roomData.RoomKey,
		RoomType:   roomData.RoomType,
		SecretFlag: roomData.SecretFlag,
	}

	chatLineDataInput := input.ChatLineDataInput{
		Cmd:          lineData.Cmd,
		LineKey:      lineData.LineKey,
		Contents:     lineData.Contents,
		SendUserHash: lineData.SendUserHash,
		SendDate:     lineData.SendDate,
	}
	return input.SendChatInput{
		EventType:    eventType,
		ChatSession:  chatSession,
		ChatRoomData: chatRoomDataInput,
		ChatLineData: chatLineDataInput,
	}
}
