package model

import "time"

type ServiceUsers struct {
	Org      string    `gorm:"column:org;varchar(30);primaryKey;comment:org"`
	UserHash string    `gorm:"column:user_hash;varchar(191);primaryKey;comment:사용자 hash 정보"`
	UserId   string    `gorm:"column:user_id;varchar(191);comment:사용자 ID"`
	UseYn    string    `gorm:"column:use_yn;varchar(3);default:Y;comment:사용여부"`
	CreateAt time.Time `gorm:"column:create_at;autoCreateTime;comment:DB 저장시간"`
}

func (ServiceUsers) TableName() string {
	return "service_users"
}

// 20260101 정리
// CREATE TABLE `service_users` (
//   `org` varchar(30) NOT NULL COMMENT 'org',
//   `user_hash` varchar(191) NOT NULL COMMENT '사용자 hash 정보 ',
//   `user_id` varchar(191) NOT NULL COMMENT '사용자 ID',
//   `use_yn` varchar(3) DEFAULT 'Y' COMMENT '사용여부',
//   `create_at` datetime(3) DEFAULT NULL COMMENT 'DB 저장시간',
//   PRIMARY KEY (`org`,`user_hash`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
