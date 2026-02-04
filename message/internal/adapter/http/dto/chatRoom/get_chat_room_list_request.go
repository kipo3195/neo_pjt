package chatRoom

type GetChatRoomListRequest struct {
	RoomType string `json:"roomType" validate:"required"`
	ReqCount int    `json:"reqCount" validate:"required"`
	Hash     string `json:"hash"`    // 기준 시간 = line key
	Filter   string `json:"filter"`  // 필터 구분 (미확인, secretFlag, ...)
	Sorting  string `json:"sorting"` // 정렬 구분 (최신, 과거, ...)
}
