package dto

type AppTokenRefreshRequest struct {
	Uuid         string `json:"uuid"`
	RefreshToken string `json:"refreshToken"`
}
