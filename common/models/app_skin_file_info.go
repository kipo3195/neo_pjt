package models

type AppSkinFileInfo struct {
	SkinType   string `gorm:"column:skin_type;primaryKey;type:varchar(20)"`             // loginImg, worksImg, color
	Device     string `gorm:"column:device;primaryKey;type:varchar(2);comment:디바이스 종류"` // 공통 C, 안드로이드 A, 아이폰 I, 웹 W
	FileName   string `gorm:"column:file_name;comment:파일 명"`
	FileHash   string `gorm:"column:file_hash;comment:파일 해시"`
	CreateDate string `gorm:"column:create_date;comment:추가시간"`
	UpdateDate string `gorm:"column:update_date;comment:업데이트 시간"`
}

func (AppSkinFileInfo) TableName() string {
	return "app_skin_file_info"
}
