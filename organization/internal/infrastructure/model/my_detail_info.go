package model

type MyDetailInfo struct {
	UserHash string `gorm:"column:user_hash"`
	DefLang  string `gorm:"column:def_lang"`
	KoLang   string `gorm:"column:ko_lang"`
	EnLang   string `gorm:"column:en_lang"`
	ZhLang   string `gorm:"column:zh_lang"`
	JpLang   string `gorm:"column:jp_lang"`
	RuLang   string `gorm:"column:ru_lang"`
	ViLang   string `gorm:"column:vi_lang"`
}
