package auth

type AppTokenValidationRequestDTO struct {
	Header AppTokenValidationRequestHeader
	Body   AppTokenValidationRequestBody
}

type AppTokenValidationRequestBody struct {
	AppToken string `json:"appToken"`
	Uuid     string `json:"uuid"`
}

type AppTokenValidationRequestHeader struct {
	ServerToken string `json:"serverToken"`
}
