package output

type ChatFileDataOutput struct {
	FileId       string `json:"fileId"`
	FileExt      string `json:"fileExt"`
	FileName     string `json:"fileName"`
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
}
