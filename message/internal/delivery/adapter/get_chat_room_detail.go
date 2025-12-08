package adapter

import (
	"message/internal/application/usecase/input"
)

func MakeGetChatRoomDetailInput(reqUserHash string, roomType string, roomKey []string) input.GetChatRoomDetailInput {

	return input.GetChatRoomDetailInput{
		ReqUserHash: reqUserHash,
		RoomType:    roomType,
		RoomKey:     roomKey,
	}

}
