package adapter

import (
	"message/internal/application/usecase/input"
)

func MakeGetChatRoomMemberReadDateInput(roomKey string, roomType string, userHash string) input.GetChatRoomMemberReadDateInput {

	return input.GetChatRoomMemberReadDateInput{
		RoomKey:  roomKey,
		RoomType: roomType,
		UserHash: userHash,
	}
}
