package model

type UserDetail struct {
	UserHash     string `gorm:"column:user_hash;primaryKey;comment:pk"`
	UserPhoneNum string `gorm:"column:user_phone_num;comment:사용자 휴대폰 번호"`
}

func (UserDetail) TableName() string {
	return "user_detail"
}

// 사용자 상세정보
