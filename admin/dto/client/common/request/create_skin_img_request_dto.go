package common

type CreateSkinImgRequestDTO struct {
	Body   CreateSkinImgRequestBody   `json:"body"`
	Header CreateSkinImgRequestHeader `json:"header"`
}

type CreateSkinImgRequestBody struct {
	SkinType string `json:"skinType"`
	File     []byte `json:"file"`
	FileSize int64  `json:"fileSize"`
	FileName string `json:"fileName"`
}

type CreateSkinImgRequestHeader struct {
}
