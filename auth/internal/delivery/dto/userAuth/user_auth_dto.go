package userAuth

type UserAuthRequest struct {
	Id   string `json:"id" validate:"required"`
	Fv   string `json:"fv" validate:"required"`
	Uuid string `json:"uuid" validate:"required"`
}

type UserAuthResponse struct {
	AccessToken     string `json:"accessToken"`
	RefreshToken    string `json:"refreshToken"`
	DeviceChallenge string `json:"deviceChallenge"`
}
