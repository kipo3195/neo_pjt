package models

// 인증 처리를 위한 정보

type AuthInfo struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

func (AuthInfo) TableName() string {
	return "auth_info"
}
