package output

type DeviceRegistOutput struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
}

func MakeDeviceRegistOutput(at string, rt string) DeviceRegistOutput {
	return DeviceRegistOutput{
		AccessToken:  at,
		RefreshToken: rt,
	}
}
