package adapter

import "message/internal/application/usecase/input"

func MakeGetChatRoomUpdateDateInput(reqUserHash string, t string, date string) input.GetChatRoomUpdateInput {

	return input.GetChatRoomUpdateInput{
		ReqUserHash: reqUserHash,
		Type:        t,
		Date:        date,
	}
}
