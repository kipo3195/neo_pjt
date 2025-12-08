package input

type GetChatRoomListInput struct {
	ReqUserHash string
	RoomType    string
	Hash        string
	ReqCount    int
	Filter      string
	Sorting     string
}
