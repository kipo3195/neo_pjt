package mapper

import (
	"message/internal/application/usecase/input"
)

func MakeDeleteChatRoomTitleInput(userHash string, org string, roomKey string, t string) input.DeleteChatRoomTitleInput {

	return input.DeleteChatRoomTitleInput{
		UserHash: userHash,
		Org:      org,
		RoomKey:  roomKey,
		Type:     t,
	}
}
