package dto

type Token struct {
	AccessToken     string `json:"accessToken"`
	RefreshToken    string `json:"refreshToken"`
	RefreshTokenExp string `json:"refreshTokenExp"`
}
