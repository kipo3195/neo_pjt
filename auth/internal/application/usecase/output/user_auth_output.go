package output

type UserAuthOutput struct {
	AccessToken     string `json:"accessToken"`
	RefreshToken    string `json:"refreshToken"`
	DeviceChallenge string `json:"deviceChallenge"`
}

func MakeUserAuthOutput(at string, rt string, dc string) UserAuthOutput {
	return UserAuthOutput{
		AccessToken:     at,
		RefreshToken:    rt,
		DeviceChallenge: dc,
	}
}
