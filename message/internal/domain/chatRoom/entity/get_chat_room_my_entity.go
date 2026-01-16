package entity

type GetChatRoomMyEntity struct {
	ReqUserHash string
	WorksCode   string
}

func MakeGetChatRoomMyEntity(reqUserHash string, worksCode string) GetChatRoomMyEntity {
	return GetChatRoomMyEntity{
		ReqUserHash: reqUserHash,
		WorksCode:   worksCode,
	}
}
