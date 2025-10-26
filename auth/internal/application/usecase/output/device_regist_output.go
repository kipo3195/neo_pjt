package output

type DeviceRegistOutput struct {
	RefreshToken    string `json:"refreshToken"`
	AccessToken     string `json:"accessToken"`
	RefreshTokenExp string `json:"refreshTokenExp"`
}
