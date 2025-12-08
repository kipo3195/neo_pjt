package entity

type GetChatRoomListEntity struct {
	ReqUserHash string
	RoomType    string
	Hash        string
	ReqCount    int
	Filter      string
	Sorting     string
}

func MakeGetChatRoomListEntity(reqUserHash string, roomType string, hash string, reqCount int, filter string, sorting string) GetChatRoomListEntity {
	return GetChatRoomListEntity{
		ReqUserHash: reqUserHash,
		RoomType:    roomType,
		Hash:        hash,
		ReqCount:    reqCount,
		Filter:      filter,
		Sorting:     sorting,
	}
}
