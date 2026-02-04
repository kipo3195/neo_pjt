package mapper

import (
	"message/internal/application/usecase/input"
)

func MakeGetChatRoomMyInput(reqUserHash string, worksCode string) input.GetChatRoomMyInput {

	return input.GetChatRoomMyInput{
		ReqUserHash: reqUserHash,
		WorksCode:   worksCode,
	}
}
