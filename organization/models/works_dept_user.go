package models

import "time"

type WorksDeptUser struct {
	Seq                  int       `gorm:"column:seq;primaryKey;comment:pk"`
	DeptCode             string    `gorm:"column:dept_code;primaryKey;comment:부서 코드"`
	DeptOrg              string    `gorm:"column:dept_org;comment:부서의 계열사 코드"`
	UserHash             string    `gorm:"column:user_hash;comment:사용자 hash 정보"`
	PositionCode         string    `gorm:"column:position_code;comment:직위 - 사원, 대리, 과장, 차장, 부장"`
	RoleCode             string    `gorm:"column:role_code;comment:직책 - 파트장, 팀장, 본부장"`
	RankNumber           string    `gorm:"column:rank_number;commnt:직급 (직위의 세부구분) - 1호봉, 2호봉, 3호봉"`
	UseYn                string    `gorm:"column:use_yn;default:Y;comment:사용 여부"`
	IsConcurrentPosition string    `gorm:"column:is_concurrent_position;default:N;comment:겸직 여부"`
	DeptCreateDate       time.Time `gorm:"column:dept_create_date;default:CURRENT_TIMESTAMP;comment:등록일"`
	UpdateHash           string    `gorm:"column:update_hash;comment:해시 정보"`
}

func (WorksDeptUser) TableName() string {
	return "works_dept_user"
}

// 부서 내 사용자
