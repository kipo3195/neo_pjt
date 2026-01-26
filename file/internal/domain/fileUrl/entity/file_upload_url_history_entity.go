package entity

type FileUploadUrlHistoryEntity struct {
	FileId     string `gorm:"column:file_id"`
	ErrorFlag  string `gorm:"column:error_flag"`
	UploadFlag string `gorm:"column:upload_flag"`
}
