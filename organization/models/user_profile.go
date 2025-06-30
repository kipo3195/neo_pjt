package models

type UserProfile struct {
	UserHash   string `gorm:"column:user_hash;primaryKey;comment:pk"`
	ProfileUrl string `gorm:"column:profile_url;type:varchar(200);comment:프로필 사진 url 정보"`
	ProfileMsg string `gorm:"column:profile_msg;type:varchar(200);comment:프로필 상태 메시지"`
}

func (UserProfile) TableName() string {
	return "user_profile"
}

// 사용자 프로필 & 상태메시지 정보
