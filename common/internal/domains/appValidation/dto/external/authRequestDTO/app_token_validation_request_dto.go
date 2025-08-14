package authRequestDTO

type AppTokenValidationRequestDTO struct {
	Header AppTokenValidationRequestHeader
	Body   AppTokenValidationRequestBody
}

type AppTokenValidationRequestBody struct {
	Token     string `json:"token"`
	TokenType string `json:"tokenType"`
	Uuid      string `json:"uuid"`
}

type AppTokenValidationRequestHeader struct {
	ServerToken string `json:"serverToken"`
}
