package fileUrl

type FileInfoDto struct {
	FileId   string `json:"fileId"`
	FileName string `json:"fileName"`
	FileSize int    `json:"fileSize"`
	FileExt  string `json:"fileExt"`
}
