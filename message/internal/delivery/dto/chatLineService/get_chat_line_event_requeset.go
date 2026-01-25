package chatLineService

type GetChatLineEventRequest struct {
	Org      string `json:"org" validate:"required"`
	RoomType string `json:"roomType" validate:"required"`
	LineKey  string `json:"lineKey" validate:"required"`
	RoomKey  string `json:"roomKey" validate:"required"`
}
