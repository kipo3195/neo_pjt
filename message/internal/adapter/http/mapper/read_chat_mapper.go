package mapper

import "message/internal/application/usecase/input"

func MakeReadChatInput(roomKey string, roomType string, userHash string, readDate string) input.ReadChatInput {
	return input.ReadChatInput{
		RoomKey:  roomKey,
		RoomType: roomType,
		UserHash: userHash,
		ReadDate: readDate,
	}
}
