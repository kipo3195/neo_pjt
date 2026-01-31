package output

type ChatFileDataOutput struct {
	FileId       string `json:"fileId"`
	FileType     string `json:"fileType"`
	FileName     string `json:"fileName"`
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
}
