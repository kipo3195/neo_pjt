package responseDTO

// 인증 결과 body
type AuthResponse struct {
	//Result       string `json:"result"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	//ConfigKey    string `json:"configKey"`
}
