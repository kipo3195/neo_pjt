package entity

type ChatLineEventEntity struct {
	EventType     string `gorm:"column:event_type" `
	Cmd           int    `gorm:"column:cmd"`
	LineKey       string `gorm:"line_key"`
	TargetLineKey string `gorm:"target_line_key"`
	Contents      string `gorm:"contents"`
	SendUserHash  string `gorm:"send_user_hash"`
	SendDate      string `gorm:"send_date"`
}
