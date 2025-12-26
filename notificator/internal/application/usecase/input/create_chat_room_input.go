package input

import "time"

type CreateChatRoomInput struct {
	CreateUserHash string    `json:"createUserHash"`
	RegDate        time.Time `json:"regDate"`
	RoomKey        string    `json:"roomKey"`
	RoomType       string    `json:"roomType"`
	Title          string    `json:"title"`
	SecretFlag     string    `json:"secretFlag"`
	Secret         string    `json:"secret"`
	Description    string    `json:"description"`
	WorksCode      string    `json:"worksCode"`
}
