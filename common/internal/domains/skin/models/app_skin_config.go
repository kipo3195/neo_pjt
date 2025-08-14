package models

type AppSkinConfig struct {
	WorksCode  string `gorm:"column:works_code;comment:works code"`
	Kind       string `gorm:"column:kind;primaryKey;comment:설정 종류"`
	Value      string `gorm:"column:value;comment:설정 값"`
	CreateDate string `gorm:"column:create_date;comment:추가시간"`
	UpdateDate string `gorm:"column:update_date;comment:업데이트 시간"`
}

func (AppSkinConfig) TableName() string {
	return "app_skin_config"
}
