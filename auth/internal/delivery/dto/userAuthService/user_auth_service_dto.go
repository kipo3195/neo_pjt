package userAuthService

type UserAuthServiceRequest struct {
	Id   string `json:"id" validate:"required"`
	Fv   string `json:"fv" validate:"required"`
	Uuid string `json:"uuid" validate:"required"`
}

type UserAuthServiceResponse struct {
	AccessToken     string `json:"accessToken"`
	RefreshToken    string `json:"refreshToken"`
	DeviceChallenge string `json:"deviceChallenge"`
}
