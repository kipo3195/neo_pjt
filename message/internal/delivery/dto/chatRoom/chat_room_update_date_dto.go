package chatRoom

type ChatRoomUpdateDateDto struct {
	RoomKey  string `json:"roomKey"`
	RoomType string `json:"roomType"`
	Detail   string `json:"detail,omitempty"`
	Line     string `json:"line,omitempty"`
	Member   string `json:"member,omitempty"`
	Owner    string `json:"owner,omitempty"`
	Title    string `json:"title,omitempty"`
}
