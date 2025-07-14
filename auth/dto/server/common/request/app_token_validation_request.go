package common

type AppTokenValidationRequestDTO struct {
	Body AppTokenValidationRequestBody
}

type AppTokenValidationRequestBody struct {
	AppToken string `json:"appToken" validate:"required"`
	Uuid     string `json:"uuid" validate:"required"`
}
