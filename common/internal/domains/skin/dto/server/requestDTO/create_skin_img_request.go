package requestDTO

import "mime/multipart"

type CreateSkinImgRequest struct {
	SkinType string                `json:"skinType"`
	File     multipart.File        `json:"file"`
	FileInfo *multipart.FileHeader `json:"fileInfo"`
}
