package fileUrl

type CreateFileUrlRequest struct {
	FileInfo  []FileInfoDto `json:"fileInfo" validate:"required"`
	Org       string        `json:"org" validate:"required"`
	EventType string        `json:"eventType" validate:"required"` // 용도 프로필, 채팅, 쪽지 ...
}
