package adapter

import (
	"notificator/internal/application/usecase/input"
	"notificator/internal/application/usecase/output"
)

func MakeChatInput(eventType string, chatSession string, roomData output.ChatRoomDataOutput, lineData output.ChatLineDataOutput) input.ChatInput {

	chatRoomDataInput := input.ChatRoomDataInput{
		RoomKey:    roomData.RoomKey,
		RoomType:   roomData.RoomType,
		SecretFlag: roomData.SecretFlag,
	}

	chatLineDataInput := input.ChatLineDataInput{
		Cmd:           lineData.Cmd,
		LineKey:       lineData.LineKey,
		TargetLineKey: lineData.TargetLineKey,
		Contents:      lineData.Contents,
		SendUserHash:  lineData.SendUserHash,
		SendDate:      lineData.SendDate,
	}
	return input.ChatInput{
		EventType:    eventType,
		ChatSession:  chatSession,
		ChatRoomData: chatRoomDataInput,
		ChatLineData: chatLineDataInput,
	}
}
