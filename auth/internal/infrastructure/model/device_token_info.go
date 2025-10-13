package model

type DeviceTokenInfo struct {
	TokenType string `gorm:"column:token_type;type:varchar(30);comment:토큰 타입"`
	TokenExp  int    `gorm:"column:token_exp;type:int(11);comment:토큰 만료 (일)"`
}

func (DeviceTokenInfo) TableName() string {
	return "device_token_info"
}
