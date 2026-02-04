package model

type OtpKeyRegistHistory struct {
	Id           string `gorm:"column:id;type:varchar(191);primaryKey;not null;comment:id"`              // UUID 또는 ID
	Uuid         string `gorm:"column:uuid;type:varchar(191);primaryKey;not null;comment:uuid"`          // 추가 PK
	OtpKey       string `gorm:"column:otp_key;type:varchar(500);not null;comment:otp 키"`                 // 암호화된 채팅 키
	Kind         string `gorm:"column:kind;type:varchar(10);primaryKey;not null;comment:키 종류"`           // 키 종류
	RegDate      string `gorm:"column:reg_date;not null;comment:등록 시간"`                                  // 발급 시간
	SvKeyVersion string `gorm:"column:sv_key_version;type:varchar(3);primaryKey;not null; comment:키 버전"` // 키 버전
}

func (OtpKeyRegistHistory) TableName() string {
	return "otp_key_regist_history"
}

// id 기반으로 등록하는 이유는 userHash가 너무 길기 때문이다 + device 등록 시점에 otp를 등록하는데 해당 시점에는 AT가 없으므로 hash를 뽑아낼 수 없다.
// 20260101 정리
