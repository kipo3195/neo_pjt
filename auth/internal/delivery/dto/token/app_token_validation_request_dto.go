package token

type AppTokenValidationRequestDTO struct {
	Body AppTokenValidationRequestBody
}

type AppTokenValidationRequestBody struct {
	AppToken string `json:"appToken" validate:"required"`
	// Uuid     string `json:"uuid" validate:"required"`
	// 20250819 수정 tokenType에 따라 appToken인지 accessToken인지 검증하는 로직
	// http://bookstack.ucware.local/books/neo-api/page/skin-config-hash-api
	Token     string `json:"token"`
	TokenType string `json:"tokenType"`
	Uuid      string `json:"uuid"`
}
