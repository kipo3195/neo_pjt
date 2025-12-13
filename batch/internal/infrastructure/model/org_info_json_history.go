package model

import "time"

type OrgInfoJsonHistory struct {
	Seq         int       `gorm:"column:seq;type:int(11);autoIncrement;primaryKey;comment:pk"`
	Org         string    `gorm:"column:org;type:varchar(30);not null;comment:org 코드"`
	FileName    string    `gorm:"column:file_name;type:varchar(50);not null;comment:생성된 json 파일명"`
	OrgInfoJson string    `gorm:"coulmn:org_info_json;type:mediumtext;comment:연동시 생성된 json"`
	CreateDate  time.Time `gorm:"column:create_date;default:CURRENT_TIMESTAMP;comment:등록일"`
}

func (OrgInfoJsonHistory) TableName() string {
	return "org_info_json_history"
}
