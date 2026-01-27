package model

type ChatFileHistory struct {
	FileId      string `gorm:"column:file_id;primaryKey;type:varchar(100);not null;comment:file id"`
	LineKey     string `gorm:"column:line_key;type:varchar(50);not null;comment:라인 키"`
	FileName    string `gorm:"column:file_name;type:varchar(200);not null;comment:file 명"`
	ReqUserHash string `gorm:"column:req_user_hash;varchar(191);not null;comment:요청 사용자 hash"`
	FileType    string `gorm:"column:file_type;varchar(10);not null;comment:파일 타입 img, file, video"`
	SendDate    string `gorm:"column:send_date;type:varchar(25);not null;comment:발송시간(서버 기준)"`
}

func (ChatFileHistory) TableName() string {
	return "chat_file_history"
}
