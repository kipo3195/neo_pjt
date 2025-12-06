package input

type CreateChatRoomInput struct {
	CreateUserHash string
	RoomKey        string
	RoomType       string
	Title          string
	SecretFlag     string
	Secret         string
	Description    string
	WorksCode      string
	Member         []CreateChatMemberInput
}
