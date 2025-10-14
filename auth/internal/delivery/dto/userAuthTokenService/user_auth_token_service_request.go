package userAuthTokenService

type AccessTokenReIssueRequest struct {
	AppToken     string `json:"appToken" validate:"required"`
	TokenType    string `json:"tokenType" validate:"required"`
	Token        string `json:"token" validate:"required"`
	Uuid         string `json:"uuid" validate:"required"`
	RefreshToken string `json:"refreshToken" validate:"required"`
}
