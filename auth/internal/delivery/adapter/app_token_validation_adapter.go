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

func MakeCheckRefreshTokenInput(userId string, uuid string, refreshToken string) input.CheckRefreshTokenInput {

	return input.CheckRefreshTokenInput{
		Uuid:         uuid,
		RefreshToken: refreshToken,
		UserId:       userId,
	}
}

func MakeReIssueAccessTokenInput(userId string, uuid string) input.ReIssueAccessTokenInput {
	return input.ReIssueAccessTokenInput{
		UserId: userId,
		Uuid:   uuid,
	}
}
