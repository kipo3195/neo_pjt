package output

type GenerateAuthTokenOutput struct {
	RefreshToken    string `json:"refreshToken"`
	AccessToken     string `json:"accessToken"`
	RefreshTokenExp string `json:"refreshTokenExp"`
}
