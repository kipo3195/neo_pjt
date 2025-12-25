package input

type UpdateChatRoomTitleInput struct {
	UserHash string
	Org      string
	RoomKey  string
	Type     string
	Title    string
}
