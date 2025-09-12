package userAuth

type UserAuthRequest struct {
	Id     string `json:"id"`
	Fv     string `json:"fv"`
	Device string `json:"device"`
}

type UserAuthResponse struct {
	AccessToken     string `json:"accessToken"`
	RefreshToken    string `json:"refreshToken"`
	DeviceChallenge string `json:"deviceChallenge"`
}
