package commonResponseDTO

type IssuedAppTokenDTO struct {
	AppToken     string `json:"appToken"`
	RefreshToken string `json:"refreshToken"`
}
