package usecases

import (
	"auth/dto"
	"auth/entities"
	"auth/repositories"
)

type authUsecase struct {
	repo repositories.AuthRepository
}

type AuthUsecase interface {
	GetAuth(dto.AuthRequest) (*entities.Auth, error)
}

func NewAuthUsecase(repo repositories.AuthRepository) AuthUsecase {
	return &authUsecase{repo: repo}
}

func (u *authUsecase) GetAuth(req dto.AuthRequest) (*entities.Auth, error) {

	auth, err := u.repo.GetAuth(req)
	if err != nil {
		return nil, err
	}

	var result string
	var accessToken string
	var refreshToken string
	var configKey string
	// 인증정보 없음.
	if auth.Id == "" {
		result = "fail"
		accessToken = ""
		refreshToken = ""
		configKey = ""
	} else {
		result = "success"
		accessToken = getAccessToken()
		refreshToken = getRefreshToken()
		configKey = getConfigkey()
	}

	response := &entities.Auth{
		Result: result, AccessToken: accessToken, RefreshToken: refreshToken, ConfigKey: configKey,
	}

	return response, err
}

func getConfigkey() string {

	return ""
}

func getAccessToken() string {

	return ""
}

func getRefreshToken() string {

	return ""
}
