package requestDTO

type AppTokenValidationRequest struct {
	AppToken string `json:"appToken" validate:"required"`
	Uuid     string `json:"uuid" validate:"required"`
}
