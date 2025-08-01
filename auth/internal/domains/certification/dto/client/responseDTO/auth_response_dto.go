package responseDTO

type AuthResponseDTO struct {
	Body AuthResponseBody
}

type AuthResponseBody struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
