package dto

type AppValidationRequest struct {
	Type   string `json:"type"`
	Domain string `json:"domain"`
	Code   string `json:"code"`
}
