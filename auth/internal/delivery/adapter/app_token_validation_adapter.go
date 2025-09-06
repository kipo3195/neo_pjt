package adapter

import (
	"auth/internal/application/usecase/input"
	"auth/internal/delivery/dto/token"
)

func MakeAppTokenValidationInput(req token.AppTokenValidationRequestDTO) input.AppTokenValidationInput {

	return input.AppTokenValidationInput{
		AppToken:  req.Body.AppToken,
		Token:     req.Body.Token,
		TokenType: req.Body.TokenType,
		Uuid:      req.Body.Uuid,
	}
}
