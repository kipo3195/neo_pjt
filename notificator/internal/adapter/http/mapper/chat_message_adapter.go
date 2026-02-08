package mapper

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
