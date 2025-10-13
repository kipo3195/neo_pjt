package output

type DeviceRegistStateOutput struct {
	AccessToken           string `json:"accessToken"`
	RefreshToken          string `json:"refreshToken"`
	RefreshTokenExp       string `json:"refreshTokenExp"`
	DeviceRegistChallenge string `json:"deviceRegistChallenge"`
}
