package model

import "time"

type IssuedAuthTokenHistory struct {
	Seq             int       `gorm:"column:seq;primaryKey;autoIncrement;comment:pk"`
	Id              string    `gorm:"column:id;type:varchar(100);comment:사용자 로그인 id"`
	Uuid            string    `gorm:"column:uuid;type:varchar(100);comment:기기 고유값"`
	AccessToken     string    `gorm:"column:access_token;type:varchar(400);comment:발급된 access 토큰 정보 JWT"`
	RefreshToken    string    `gorm:"column:refresh_token;type:varchar(400);comment:발급된 refresh 토큰 정보 JWT"`
	RefreshTokenExp string    `gorm:"column:refresh_token_exp;type:varchar(20);comment:발급된 refresh 토큰의 만료날짜"`
	CreateAt        time.Time `gorm:"column:create_at;autoCreateTime;comment:DB 저장시간"`
	UpdateAt        time.Time `gorm:"column:update_at;autoUpdateTime;comment:DB 업데이트 시간"`
}

func (IssuedAuthTokenHistory) TableName() string {
	return "issued_auth_token_history"
}

// 추후 uuid, rt 등에 index 추가 필요
