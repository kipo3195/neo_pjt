package input

type CreateSkinImgInput struct {
	SkinType string
	File     []byte
	FileSize int64
	FileName string
}

func MakeCreateSkinImgInput(skinType string, file []byte, fileSize int64, fileName string) CreateSkinImgInput {
	return CreateSkinImgInput{
		SkinType: skinType,
		File:     file,
		FileSize: fileSize,
		FileName: fileName,
	}
}
