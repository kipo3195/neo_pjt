package entity

type ChatLineEventEntity struct {
	EventType     string `gorm:"column:event_type" `
	Cmd           int    `gorm:"column:cmd"`
	LineKey       string `gorm:"line_key"`
	TargetLineKey string `gorm:"target_line_key"`
	Contents      string `gorm:"contents"`
	SendUserHash  string `gorm:"send_user_hash"`
	SendDate      string `gorm:"send_date"`
	FileId        string `gorm:"file_id"`
	FileName      string `gorm:"file_name"`
	FileType      string `gorm:"file_type"`
}
