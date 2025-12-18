package entity

type UserDetailEntity struct {
	Org          string `gorm:"column:org" json:"org"` // 추가
	UserHash     string `gorm:"column:user_hash" json:"userHash"`
	UserId       string `gorm:"column:user_id" json:"userId"`
	UserPhoneNum string `gorm:"column:user_phone_num" json:"userPhoneNum"`
	UserEmail    string `gorm:"column:user_email" json:"userEmail"`
}
