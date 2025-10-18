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

func MakeCheckRefreshTokenInput(userId string, uuid string, refreshToken string, withoutId bool) input.RefreshTokenCheckInput {

	return input.RefreshTokenCheckInput{
		Uuid:         uuid,
		RefreshToken: refreshToken,
		UserId:       userId,
		WithoutId:    withoutId,
	}
}

func MakeReIssueAccessTokenInput(userId string, uuid string) input.ReIssueAccessTokenInput {
	return input.ReIssueAccessTokenInput{
		UserId: userId,
		Uuid:   uuid,
	}
}

func MakeReIssueAccessTokenSavedInput(userId string, uuid string, rt string, at string) input.ReIssueAccessTokenSavedInput {
	return input.ReIssueAccessTokenSavedInput{
		UserId: userId,
		Uuid:   uuid,
		Rt:     rt,
		At:     at,
	}
}
