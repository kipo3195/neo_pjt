package model

type ChatLineEvent struct {
	Seq          int    `gorm:"column:seq;primaryKey;autoIncrement;not null;comment:PK"`
	EventType    string `gorm:"column:event_type;type:varchar(3);not null;comment:이벤트 타입 C, U, D"`
	Cmd          int    `gorm:"column:cmd;type:int(11);comment:cmd"`
	LineKey      string `gorm:"column:line_key;type:varchar(50);not null;comment:라인 키"`
	Contents     string `gorm:"column:contents;type:mediumtext;comment:채팅 내용"`
	SendUserHash string `gorm:"column:send_user_hash;not null;comment:발신자"`
	SendDate     string `gorm:"column:send_date;type:varchar(20);not null;comment:발송시간(서버 기준)"`
}

func (ChatLineEvent) TableName() string {
	return "chat_line_event"
}
