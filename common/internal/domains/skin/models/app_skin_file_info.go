package models

import "time"

type AppSkinFileInfo struct {
	SkinType   string    `gorm:"column:skin_type;primaryKey;type:varchar(20);comment:loginImg, worksImg 같은 스킨의 타입"` // loginImg, worksImg
	FileUrl    string    `gorm:"column:file_url;comment:파일 다운로드 url"`
	FileName   string    `gorm:"column:file_name;comment:파일 명"`
	FilePath   string    `grom:"column:file_path;comment:파일 저장 경로"`
	FileHash   string    `gorm:"column:file_hash;comment:파일 해시"` // 조회 조건
	CreateDate time.Time `gorm:"column:create_date;autoCreateTime;comment:추가시간"`
	UpdateDate time.Time `gorm:"column:update_date;autoUpdateTime;comment:업데이트 시간"`
}

func (AppSkinFileInfo) TableName() string {
	return "app_skin_file_info"
}
