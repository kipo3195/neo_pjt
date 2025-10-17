package deviceAuthService

type DeviceRefreshResponse struct {
	AccessToken     string `json:"accessToken"`
	RefreshToken    string `json:"refreshToken"`
	RefreshTokenExp string `json:"refreshTokenExp"`
}
