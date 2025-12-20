package entity

type ProfileInfoEntity struct {
	Org            string `gorm:"column:org"`
	UserHash       string `gorm:"column:user_hash"`
	SaveName       string `gorm:"column:save_name"`
	ProfileVersion int64
	ProfileMsg     string `gorm:"column:profile_msg"`
}
