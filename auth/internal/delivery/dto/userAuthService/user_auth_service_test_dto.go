package userAuthService

type UserAuthServiceTestRequest struct {
	UserId string `json:"userId"`
	Uuid   string `json:"uuid"`
}

type UserAuthServiceTestResponse struct {
	AccessToken     string `json:"accessToken"`
	RefreshToken    string `json:"refreshToken"`
	RefreshTokenExp string `json:"refreshTokenExp"`
}
