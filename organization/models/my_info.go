package models

type MyInfo struct {
	UserHash     string `gorm:"column:user_hash"`
	UserPhoneNum string `gorm:"column:user_phone_num"`
	KrLang       string `gorm:"column:kr_lang"`
	EnLang       string `gorm:"column:en_lang"`
	CnLang       string `gorm:"column:cn_lang"`
	JpLang       string `gorm:"column:jp_lang"`
}
