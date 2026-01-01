package model

import "time"

type ServiceUsers struct {
	Org      string    `gorm:"column:org;varchar(30);primaryKey;comment:works code"`
	UserId   string    `gorm:"column:user_id;varchar(191);primaryKey;comment:사용자 ID"`
	UserHash string    `gorm:"column:user_hash;varchar(191);comment:사용자 hash 정보"`
	UseYn    string    `gorm:"column:use_yn;varchar(3);default:Y;comment:사용여부"`
	CreateAt time.Time `gorm:"column:create_at;autoCreateTime;comment:DB 저장시간"`
}

func (ServiceUsers) TableName() string {
	return "service_users"
}

// 서비스 등록 사용자 (공통) 20260101 기준으로 작성.
// user_id는 변할 수 있지만, user_hash는 변할 수 없으므로 pk는 org와 userId로 처리한다.
// 단, 필요에 의해 org, user_hash를 pk로 두고 `연동 계정` 이라는 별도의 테이블을 생성해서 user_hash로 매핑 하는 방안도 사용할 수 있다.
// id, hash의 길이가 191인 것은 gorm의 autoMigrate에 의해 테이블이 생성될때 string에 매핑되는 고유 길이다.
