package chatRoom

type CreateChatRoomRequest struct {
	RoomKey     string           `json:"roomKey" validate:"required"`
	Type        string           `json:"roomType" validate:"required"`
	Title       string           `json:"title" validate:"required"`
	SecretFlag  string           `json:"secretFlag" validate:"required"`
	Secret      string           `json:"secret"`
	Description string           `json:"description" validate:"required"`
	WorksCode   string           `json:"worksCode" validate:"required"`
	Member      []ChatRoomMember `json:"member" validate:"required"`
}
