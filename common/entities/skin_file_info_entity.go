package entities

type SkinFileInfoEntity struct {
	FileHash string `json:"fileHash"`
	FileName string `json:"fileName"`
	SkinType string `json:"skinType"`
	Device   string `json:"device"`
}
