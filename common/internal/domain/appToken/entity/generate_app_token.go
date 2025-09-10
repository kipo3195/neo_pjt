package entity

type GenerateAppToken struct {
	AppToken     string `json:"appToken"`
	Uuid         string `json:"uuid"`
	RefreshToken string `json:"refreshToken"`
}
