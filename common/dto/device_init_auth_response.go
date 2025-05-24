package dto

// auth에서 전달받은 토큰 처리용  dto

type DeviceInitAuthResponse struct {
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
}
