package adapter

import "admin/internal/application/usecase/input"

func MakeCreateSkinImgInput(skinType string, file []byte, fileSize int64, fileName string) input.CreateSkinImgInput {
	return input.CreateSkinImgInput{
		SkinType: skinType,
		File:     file,
		FileSize: fileSize,
		FileName: fileName,
	}
}
