package chatRoom

type GetChatRoomUpdateDateRequest struct {
	Type string `json:"type"` // A 전체, S 일부
	Date string `json:"date"` // type이 S인 경우 기준 데이터
}
