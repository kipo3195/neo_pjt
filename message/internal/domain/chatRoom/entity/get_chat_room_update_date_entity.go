package entity

type GetChatRoomUpdateDateEntity struct {
	ReqUserHash string
	Type        string
	Date        string
}

func MakeGetChatRoomUpdateDateEntity(reqUserHash string, t string, date string) GetChatRoomUpdateDateEntity {
	return GetChatRoomUpdateDateEntity{
		ReqUserHash: reqUserHash,
		Type:        t,
		Date:        date,
	}
}
