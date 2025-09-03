package appValidation

type IssuedAppTokenDTO struct {
	AppToken     string `json:"appToken"`
	RefreshToken string `json:"refreshToken"`
}
