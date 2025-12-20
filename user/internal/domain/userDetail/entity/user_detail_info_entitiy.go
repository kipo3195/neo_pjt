package entity

type UserDetailInfoEntity struct {
	Org           string `gorm:"column:org"`
	UserHash      string `gorm:"column:user_hash"`
	UserPhoneNum  string `gorm:"column:user_phone_num"`
	UserEmail     string `gorm:"column:user_email"`
	DetailVersion int64
}
