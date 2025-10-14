package adapter

import (
	"auth/internal/application/usecase/input"
)

func MakeAppTokenValidationInput(appToken string, token string, tokenType string, uuid string) input.AppTokenValidationInput {

	return input.AppTokenValidationInput{
		AppToken:  appToken,
		Token:     token,
		TokenType: tokenType,
		Uuid:      uuid,
	}
}
