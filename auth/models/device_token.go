package models

type DeviceToken struct {
	Seq   int    `gorm:"primaryKey;autoIncrement" json:"seq"`
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
}

func (DeviceToken) TableName() string {
	return "device_token"
}
