package dto

// 인증 결과 body
type AuthResponse struct {
	Result       string
	AccessToken  string
	RefreshToken string
	ConfigKey    string
}
