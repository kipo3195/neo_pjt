package entities

type AppTokenEntity struct {
	Uuid         string `json:"uuid"`
	AppToken     string `json:"appToken"`
	RefreshToken string `json:"refreshToken"`
}
