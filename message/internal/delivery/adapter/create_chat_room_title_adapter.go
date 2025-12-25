package adapter

import (
	"message/internal/application/usecase/input"
)

func MakeUpdateChatRoomTitleInput(userHash string, org string, roomKey string, t string, title string) input.UpdateChatRoomTitleInput {

	return input.UpdateChatRoomTitleInput{
		UserHash: userHash,
		Org:      org,
		RoomKey:  roomKey,
		Type:     t,
		Title:    title,
	}
}
