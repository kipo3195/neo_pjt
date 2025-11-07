package model

import "time"

type ProfileImgInfo struct {
	UserHash            string    `gorm:"column:user_hash;type:varchar(100);primaryKey;	comment:사용자 hash 정보"`
	UserId              string    `gorm:"column:user_id;type:varchar(100);comment:사용자 계정"`
	ProfileImgHash      string    `gorm:"column:img_hash;type:varchar(100);primaryKey;comment:이미지 hash 정보"`
	ProfileImgSavedName string    `gorm:"column:save_name;type:varchar(100);comment:이미지 이미지 명칭 정보"`
	ProfileImgSavedPath string    `gorm:"column:save_path;type:varchar(300);comment:이미지 저장 경로 정보"`
	ProfileImgSize      int64     `gorm:"column:size;type:int(11);comment:이미지 사이즈"`
	CreateAt            time.Time `gorm:"column:create_at;autoCreateTime;comment:DB 저장시간"`
	UseYn               string    `gorm:"column:use_yn;type:varchar(1);default:Y;comment:사용 유무"`
	UpdateAt            time.Time `gorm:"column:update_at;autoUpdateTime;comment:DB 삭제 시간"`
}

func (ProfileImgInfo) TableName() string {
	return "profile_img_info"
}

// id, profileImgHash를 복합키로 잡은 이유?
// 프로필 이미지는 여러개 등록 가능하게 처리여 추후 멀티 프로필기능 추가시 대응..
// 멀티 프로필 노출 대상 table을 하나 생성해서 join 해서 사용할 수 있도록...

// 그리고 보관기간에 따른 데이터 삭제시 처리 방안도 생각해볼것. DB 테이블에서 삭제. 현재 테이블을 update방향으로 사용하지 않고 insert만 하므로
