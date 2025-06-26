package dto

type SvAppTokenValidationRequest struct {
	AppToken string `json:"appToken" validate:"required"`
	Uuid     string `json:"uuid" validate:"required"`
}
