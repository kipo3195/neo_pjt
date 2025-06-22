package models

type AppSkinConfig struct {
	WorksCode  string `gorm:"column:works_code"`
	Device     string `gorm:"column:device;primaryKey;comment:디바이스 종류"`
	Version    string `gorm:"column:version;comment:버전 정보"`
	CreateDate string `gorm:"column:create_date;comment:추가시간"`
	UpdateDate string `gorm:"column:update_date;comment:업데이트 시간"`
}

func (AppSkinConfig) TableName() string {
	return "app_skin_config"
}
