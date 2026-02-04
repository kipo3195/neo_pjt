package mapper

import "message/internal/application/usecase/input"

func MakeGetChatLineEventInput(reqUserHash string, org string, roomType string, roomKey string, lineKey string) input.GetChatLineEventInput {
	return input.GetChatLineEventInput{
		ReqUserHash: reqUserHash,
		Org:         org,
		RoomType:    roomType,
		RoomKey:     roomKey,
		LineKey:     lineKey,
	}
}
