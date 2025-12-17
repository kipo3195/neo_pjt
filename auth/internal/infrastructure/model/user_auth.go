package model

import "time"

type UserAuth struct {
	Id       string    `gorm:"column:id;type:varchar(50);primaryKey;comment:pk"`
	Salt     string    `gorm:"column:salt;type:varchar(50);comment:salt 정보"`
	AuthHash string    `gorm:"column:auth_hash;type:varchar(200);comment:인증 hash"`
	UserHash string    `gorm:"column:user_hash;type:varchar(32);comment:사용자 hash"`
	CreateAt time.Time `gorm:"column:create_at;autoCreateTime;comment:DB 저장시간"`
	UpdateAt time.Time `gorm:"column:update_at;autoCreateTime;comment:DB 저장시간"`
}

func (UserAuth) TableName() string {
	return "user_auth"
}
