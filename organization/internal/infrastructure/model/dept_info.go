package model

type DeptInfo struct {
	UserHash    string `gorm:"column:user_hash"`
	DeptOrg     string `gorm:"column:dept_org"`
	DeptCode    string `gorm:"column:dept_code"`
	DefLang     string `gorm:"column:def_lang"`
	KoLang      string `gorm:"column:ko_lang"`
	EnLang      string `gorm:"column:en_lang"`
	ZhLang      string `gorm:"column:zh_lang"`
	JpLang      string `gorm:"column:jp_lang"`
	RuLang      string `gorm:"column:ru_lang"`
	ViLang      string `gorm:"column:vi_lang"`
	Header      string `gorm:"column:header"`
	Description string `grom:"column:description"`

	// 직책(Role) 다국어
	RKrLang string `gorm:"column:r_kr_lang"`
	REnLang string `gorm:"column:r_en_lang"`
	RZhLang string `gorm:"column:r_zh_lang"`
	RJpLang string `gorm:"column:r_jp_lang"`

	// 직위(Position) 다국어
	PKrLang string `gorm:"column:p_kr_lang"`
	PEnLang string `gorm:"column:p_en_lang"`
	PZhLang string `gorm:"column:p_zh_lang"`
	PJpLang string `gorm:"column:p_jp_lang"`
}
