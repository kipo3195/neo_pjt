package common

type SkinFileInfo struct {
	FileHash string `json:"fileHash"`
	FileName string `json:"fileName"`
	SkinType string `json:"skinType"`
	FilePath string `json:"filePath"`
	FileUrl  string `json:"fileUrl"`
}
