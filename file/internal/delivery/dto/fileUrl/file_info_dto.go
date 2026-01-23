package fileUrl

type FileInfoDto struct {
	FileName string `json:"fileName"`
	FileSize int    `json:"fileSize"`
	FileExt  string `json:"fileExt"`
}
