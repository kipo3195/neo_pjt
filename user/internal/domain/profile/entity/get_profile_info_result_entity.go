package entity

type GetProfileInfoResultEntity struct {
	UserHash       string `gorm:"column:user_hash"`
	ProfileMsg     string `gorm:"column:profile_msg"`
	ProfileImgHash string `gorm:"column:img_hash"`
}
