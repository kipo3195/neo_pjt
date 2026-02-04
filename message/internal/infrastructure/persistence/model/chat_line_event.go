package model

type ChatLineEvent struct {
	Seq           int    `gorm:"column:seq;primaryKey;autoIncrement;not null;comment:PK"`
	RoomKey       string `gorm:"column:room_key;type:varchar(50);not null;comment:룸 키"`
	EventType     string `gorm:"column:event_type;type:varchar(3);not null;comment:이벤트 타입 C, U, D"`
	Cmd           int    `gorm:"column:cmd;type:int(11);comment:cmd"`
	LineKey       string `gorm:"column:line_key;type:varchar(50);not null;comment:라인 키"`
	TargetLineKey string `gorm:"column:target_line_key;type:varchar(50);comment:타겟 라인키 (이벤트 타입이 U, D일 경우)"`
	Contents      string `gorm:"column:contents;type:mediumtext;comment:채팅 내용"`
	SendUserHash  string `gorm:"column:send_user_hash;varchar(191);not null;comment:발신자 hash"`
	SendDate      string `gorm:"column:send_date;type:varchar(25);not null;comment:발송시간(서버 기준)"`
}

func (ChatLineEvent) TableName() string {
	return "chat_line_event"
}

// 20260101 정리
