package requestDTO

type AppTokenValidationRequestDTO struct {
	AppToken string `json:"appToken" validate:"required"`
	Uuid     string `json:"uuid" validate:"required"`
}
