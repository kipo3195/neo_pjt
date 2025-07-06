package auth

// auth에서 전달받은 토큰 처리용  dto

type DeviceInitAuthResponse struct {
	AppToken     string `json:"appToken"`
	RefreshToken string `json:"refreshToken"`
}
