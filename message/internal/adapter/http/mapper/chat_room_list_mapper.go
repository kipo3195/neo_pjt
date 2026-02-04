package mapper

import "message/internal/application/usecase/input"

func MakeGetChatRoomListInput(reqUserHash string, roomType string, hash string, reqCount int, filter string, sorting string) input.GetChatRoomListInput {
	return input.GetChatRoomListInput{
		ReqUserHash: reqUserHash,
		RoomType:    roomType,
		Hash:        hash,
		ReqCount:    reqCount,
		Filter:      filter,
		Sorting:     sorting,
	}
}
