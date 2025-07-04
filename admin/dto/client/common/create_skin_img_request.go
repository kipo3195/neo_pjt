package common

import "mime/multipart"

type CreateSkinImgRequest struct {
	File     multipart.File        `json:"file"`
	FileInfo *multipart.FileHeader `json:"fileInfo"`
	SkinType string                `json:"skinType"`
}
