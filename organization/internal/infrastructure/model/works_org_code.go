package model

import "time"

type WorksOrgCode struct {
	OrgCode    string    `gorm:"column:org_code;type:varchar(30);primaryKey;comment:org code"`
	CreateDate time.Time `gorm:"column:create_date;default:CURRENT_TIMESTAMP;comment:등록일"`
}

func (WorksOrgCode) TableName() string {
	return "works_org_code"
}

// org 서비스의 org code의 기준이 되는 테이블 - 관리자 기능으로 등록처리 하기, 내 부서 정보 조회때도 해당 테이블의 정보를 바탕으로 view 할 수 있는 json 생성
