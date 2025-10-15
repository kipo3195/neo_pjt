package token

type AccessTokenReIssueRequest struct {
	AppToken     string `json:"appToken" validate:"required"`
	Uuid         string `json:"uuid" validate:"required"`
	RefreshToken string `json:"refreshToken" validate:"required"`
}
