package output

type DeviceRegistOutput struct {
	RefreshToken    string `json:"refreshToken"`
	AccessToken     string `json:"accessToken"`
	RefreshTokenExp string `json:"refreshTokenExp"`
}

func MakeDeviceRegistOutput(at string, rt string, rtExp string) DeviceRegistOutput {
	return DeviceRegistOutput{
		AccessToken:     at,
		RefreshToken:    rt,
		RefreshTokenExp: rtExp,
	}
}
