package model

import "time"

type DeviceRegistHistory struct {
	Id        string    `gorm:"column:id;type:varchar(100);primaryKey;comment:사용자 계정"`
	Uuid      string    `gorm:"column:uuid;type:varchar(100);primaryKey;comment:기기 고유값"`
	ModelName string    `gorm:"column:model_name;type:varchar(100);comment:기기 모델명"`
	Version   string    `gorm:"column:version;type:varchar(100);comment:기기 버전"`
	CreateAt  time.Time `gorm:"column:create_at;autoCreateTime;comment:DB 저장시간"`
	UseYn     string    `gorm:"column:use_yn;type:varchar(1);default:Y;comment:사용유무"`
}

func (DeviceRegistHistory) TableName() string {
	return "device_regist_history"
}
