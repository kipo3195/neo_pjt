package dto

// 인증 요청 body
type AuthRequest struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}
