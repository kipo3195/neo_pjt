package models

type UserGrade struct {
	UserHash  string `gorm:"column:user_hash;primaryKey;comment:'pk'"`
	UserGrade int    `gorm:"column:user_grade;comment:'사용자 등급'"`
}

func (UserGrade) TableName() string {
	return "user_grade"
}

// 사용자 등급 정보
