package adapter

import (
	"notificator/internal/application/usecase/input"
)

func MakeChatMessageInput(eventType string, chatSession string, roomData input.ChatRoomDataInput, lineData input.ChatLineDataInput) input.ChatMessageInput {
	return input.ChatMessageInput{
		EventType:    eventType,
		ChatSession:  chatSession,
		ChatRoomData: roomData,
		ChatLineData: lineData,
	}
}

// func MakeChatMessageOutput(en entity.ChatMessage) output.ChatMessageOutput {

// 	chatRoomData := output.ChatRoomDataOutput{
// 		RoomType:   en.ChatRoomData.RoomType,
// 		RoomKey:    en.ChatRoomData.RoomKey,
// 		SecretFlag: en.ChatRoomData.SecretFlag,
// 	}

// 	chatLineData := output.ChatLineDataOutput{
// 		Cmd:           en.ChatLineData.Cmd,
// 		Contents:      en.ChatLineData.Contents,
// 		LineKey:       en.ChatLineData.LineKey,
// 		TargetLineKey: en.ChatLineData.TargetLineKey,
// 		SendUserHash:  en.ChatLineData.SendUserHash,
// 		SendDate:      en.ChatLineData.SendDate,
// 	}

// 	return output.ChatMessageOutput{
// 		ChatRoomData: chatRoomData,
// 		ChatLineData: chatLineData,
// 		EventType:    en.EventType,
// 		ChatSession:  en.ChatSession,
// 	}
// }
