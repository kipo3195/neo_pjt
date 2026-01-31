package chatService

type ChatFileData struct {
	FileId       string `json:"fileId"`
	FileName     string `json:"fileName"`
	FileType     string `json:"fileType"`
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
}
