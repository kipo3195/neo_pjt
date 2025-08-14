package entities

type SkinFileInfoEntity struct {
	FileHash string `json:"fileHash"`
	SkinType string `json:"skinType"`
	FilePath string `json:"filePath"`
}
