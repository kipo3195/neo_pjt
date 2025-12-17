package model

import "time"

type ServiceUsers struct {
	Org      string    `gorm:"column:org;varchar(30);primaryKey;comment:works code"`
	UserId   string    `gorm:"column:user_id;varchar(200);primaryKey;comment:사용자 ID"`
	UserHash string    `gorm:"column:user_hash;varchar(100);comment:사용자 hash 정보 "`
	UseYn    string    `gorm:"column:use_yn;varchar(3);default:Y;comment:사용여부"`
	CreateAt time.Time `gorm:"column:create_at;autoCreateTime;comment:DB 저장시간"`
}

func (ServiceUsers) TableName() string {
	return "service_users"
}
