package auth

type AppTokenValidationRequest struct {
	AppToken string `json:"appToken"`
	Uuid     string `json:"uuid"`
}
