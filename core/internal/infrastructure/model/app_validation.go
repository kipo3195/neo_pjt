package model

// 인증 처리를 위한 정보

type AppValidation struct {
	AppHash    string `gorm:"column:app_hash"`
	DeviceKind string `gorm:"column:device_kind"`
	UdtDate    string `gorm:"column:udt_date"`
	Version    string `gorm:"column:version"`
}

func (AppValidation) TableName() string {
	return "app_validation"
}

// CREATE TABLE `app_validation` (
//   `version_id` varchar(200) NOT NULL COMMENT '앱 해시',
//   `device_kind` varchar(10) NOT NULL COMMENT '디바이스 종류',
//   `udt_date` datetime DEFAULT current_timestamp() COMMENT '앱 버전 업데이트 시간',
//   `version` varchar(10) NOT NULL COMMENT '버전',
//   PRIMARY KEY (`version_id`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='공통 앱의 버전 관리'
