package model

import "time"

type WorksDeptUser struct {
	DeptOrg              string    `gorm:"column:dept_org;type:varchar(30);comment:부서의 계열사 코드"`
	DeptCode             string    `gorm:"column:dept_code;type:varchar(30);primaryKey;comment:부서 코드"`
	UserHash             string    `gorm:"column:user_hash;type:varchar(191);primaryKey;comment:사용자 hash 정보"`
	PositionCode         string    `gorm:"column:position_code;comment:직위 - 사원, 대리, 과장, 차장, 부장"`
	RoleCode             string    `gorm:"column:role_code;type:varchar(30);comment:직책 - 파트장, 팀장, 본부장"`
	RankNumber           string    `gorm:"column:rank_number;type:varchar(30);commnt:직급 (직위의 세부구분) - 1호봉, 2호봉, 3호봉"`
	UseYn                string    `gorm:"column:use_yn;type:varchar(3);default:Y;comment:사용 여부"`
	IsConcurrentPosition string    `gorm:"column:is_concurrent_position;type:varchar(3);default:N;comment:겸직 여부"`
	DeptCreateDate       time.Time `gorm:"column:dept_create_date;default:CURRENT_TIMESTAMP;comment:등록일"`
	UpdateHash           string    `gorm:"column:update_hash;comment:해시 정보"`
}

func (WorksDeptUser) TableName() string {
	return "works_dept_user"
}

// 부서 내 사용자
