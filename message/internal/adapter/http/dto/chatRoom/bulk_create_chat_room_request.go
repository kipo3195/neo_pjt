package chatRoom

type BulkCreateChatRoomRequest struct {
	MakeCount int `json:"makeCount" validate:"required,min=1"`
}
