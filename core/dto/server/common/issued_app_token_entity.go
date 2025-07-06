package common

type IssuedAppToken struct {
	AppToken     string `json:"appToken"`
	RefreshToken string `json:"refreshToken"`
}
