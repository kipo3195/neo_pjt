package model

type UserDetail struct {
	UserHash     string `gorm:"column:user_hash;type:varchar(100);primaryKey;comment:pk"`
	UserPhoneNum string `gorm:"column:user_phone_num;type:varchar(30);'comment:사용자 휴대폰 번호"`
	UserEmail    string `gorm:"column:user_email;type:varchar(200);comment:사용자 이메일"`
}

func (UserDetail) TableName() string {
	return "user_detail"
}

// 사용자 상세정보
