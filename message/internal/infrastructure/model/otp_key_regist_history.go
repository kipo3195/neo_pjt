package model

type OtpKeyRegistHistory struct {
	Id           string `gorm:"column:id;primaryKey;not null;comment:id"`                                // UUID 또는 ID
	Uuid         string `gorm:"column:uuid;primaryKey;not null;comment:uuid"`                            // 추가 PK
	OtpKey       string `gorm:"column:otp_key;not null;comment:otp 키"`                                   // 암호화된 채팅 키
	Kind         string `gorm:"column:kind;type:varchar(10);primaryKey;not null;comment:키 종류"`           // 키 종류
	RegDate      string `gorm:"column:reg_date;not null;comment:등록 시간"`                                  // 발급 시간
	SvKeyVersion string `gorm:"column:sv_key_version;type:varchar(3);primaryKey;not null; comment:키 버전"` // 키 버전
}

func (OtpKeyRegistHistory) TableName() string {
	return "otp_key_regist_history"
}
