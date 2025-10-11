package entity

type CreateSkinImgEntity struct {
	SkinType string
	File     []byte
	FileSize int64
	FileName string
}

func MakeCreateSkinImgEntity(skinType string, file []byte, fileSize int64, fileName string) CreateSkinImgEntity {
	return CreateSkinImgEntity{
		SkinType: skinType,
		File:     file,
		FileSize: fileSize,
		FileName: fileName,
	}
}
