package model

type UserDetail struct {
	Seq          int64  `gorm:"column:seq;primaryKey;autoIncrement;comment:pk"`
	Org          string `gorm:"column:org;type:longtext;not null;comment:org code"`
	UserId       string `gorm:"column:user_id;type:varchar(100);not null;comment:사용자 id"`
	UserHash     string `gorm:"column:user_hash;type:varchar(100);comment:사용자 hash"`
	UserPhoneNum string `gorm:"column:user_phone_num;type:varchar(30);comment:사용자 휴대폰 번호"`
	UserEmail    string `gorm:"column:user_email;type:varchar(200);comment:사용자 이메일"`
}

// TableName 명시 (복수형 방지)
func (UserDetail) TableName() string {
	return "user_detail"
}
