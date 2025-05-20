package dto

type ErrorResponse struct {
	Code    string `json:"code"`    // Error code
	Message string `json:"message"` // Error message
}
