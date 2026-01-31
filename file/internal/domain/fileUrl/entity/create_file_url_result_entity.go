package entity

type CreateFileUrlResultEntity struct {
	FileId     string `json:"fileId"`
	FileName   string `json:"fileName"`
	FileType   string `json:"fileType"`
	CreatedUrl string `json:"createdUrl"`
}
