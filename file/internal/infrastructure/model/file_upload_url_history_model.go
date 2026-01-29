package model

import "time"

type FileUploadUrlHistory struct {
	FileId        string    `gorm:"column:file_id;primaryKey;type:varchar(100);not null;comment:file id"`
	TransactionId string    `gorm:"column:t_id;type:varchar(50);not null;comment:transaction id"`
	FileName      string    `gorm:"column:file_name;type:varchar(200);not null;comment:file 명"`
	ReqUserHash   string    `gorm:"column:req_user_hash;varchar(191);not null;comment:요청 사용자 hash"`
	UploadUrl     string    `gorm:"column:upload_url;varchar(400);not null;comment:업로드 url"`
	CreateDate    time.Time `gorm:"column:create_date;not null;autoCreateTime;comment:등록 일"`
	UpdateDate    time.Time `gorm:"column:update_date;not null;autoCreateTime;autoUpdateTime;comment:수정 일"`
	UploadFlag    string    `gorm:"column:upload_flag;type:varchar(1);default:N;comment: url 발급 : N, 업로드 완료 Y"`
	SendFlag      string    `gorm:"column:send_flag;type:varchar(1);default:N;comment: message 서비스를 통한 발송완료시 Y"`
	ErrorFlag     string    `gorm:"column:error_flag;type:varchar(1);default:N;comment: url 발급 + 업로드 중: N, batch 서비스 업로드 실패 감지 : Y"`
}

func (FileUploadUrlHistory) TableName() string {
	return "file_upload_url_history"
}
