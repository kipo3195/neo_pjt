package output

type DeviceRegistStateOutput struct {
	AccessToken           string `json:"accessToken"`
	RefreshToken          string `json:"refreshToken"`
	DeviceRegistChallenge string `json:"deviceRegistChallenge"`
}
