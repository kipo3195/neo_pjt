package entity

type SendFileInfo struct {
	FileId  string `gorm:"column:file_id"`
	LineKey string `gorm:"column:line_key"`
}
