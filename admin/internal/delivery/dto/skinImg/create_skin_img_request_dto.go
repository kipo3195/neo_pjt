package skinImg

type CreateSkinImgRequest struct {
	SkinType string `json:"skinType"`
	File     []byte `json:"file"`
	FileSize int64  `json:"fileSize"`
	FileName string `json:"fileName"`
}
