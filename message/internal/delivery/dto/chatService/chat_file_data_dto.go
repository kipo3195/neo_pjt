package chatService

type ChatFileData struct {
	FileId       string `json:"fileId"`
	FileExt      string `json:"fileExt"`
	FileName     string `json:"fileName"`
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
}
