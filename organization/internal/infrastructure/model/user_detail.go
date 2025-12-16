package model

type UserDetail struct {
	UserHash     string `gorm:"column:user_hash;primaryKey;type:varchar(191)" json:"user_hash"`
	UserPhoneNum string `gorm:"column:user_phone_num;type:longtext" json:"user_phone_num"`
	UserEmail    string `gorm:"column:user_email;type:varchar(300)" json:"user_email"`
	UpdateHash   string `gorm:"column:update_hash;type:varchar(30);comment:해시 정보"`
}

func (UserDetail) TableName() string {
	return "user_detail"
}
