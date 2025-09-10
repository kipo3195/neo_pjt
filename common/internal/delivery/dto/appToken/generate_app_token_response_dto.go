package appToken

type GenerateAppTokenResponseDTO struct {
	Body GenerateAppTokenResponseBody
}

type GenerateAppTokenResponseBody struct {
	AppToken     string `json:"appToken"`
	Uuid         string `json:"uuid"`
	RefreshToken string `json:"refreshToken"`
}
