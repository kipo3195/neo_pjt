package chatRoom

type ChatRoomDetail struct {
	RoomKey     string `json:"roomKey"`
	Title       string `json:"title"`
	SecretFlag  string `json:"secretFlag"`
	Secret      string `json:"secret"`
	Description string `json:"description"`
	State       string `json:"state"`
	WorksCode   string `json:"worksCode"`
	CreateDate  string `json:"createDate"`
	CreateUser  string `json:"createUser"`
	Hash        string `json:"hash"`
}
