package input

type ChatReadMessageInput struct {
	RoomKey    string `json:"roomKey"`
	RoomType   string `json:"roomType"`
	MemberHash string `json:"memberHash"`
	ReadDate   string `json:"readDate"`
}
