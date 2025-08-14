package entities

type IssuedAppToken struct {
	AppToken     string `json:"AppToken"`
	RefreshToken string `json:"refreshToken"`
}
