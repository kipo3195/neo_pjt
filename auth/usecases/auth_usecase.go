package usecases

import (
	"auth/config"
	"auth/dto"
	"auth/entities"
	"auth/repositories"
)

type authUsecase struct {
	repo   repositories.AuthRepository
	jwtCfg *config.JWTConfig
}

type AuthUsecase interface {
	GetAuth(dto.AuthRequest) (*entities.Auth, error)
}

func NewAuthUsecase(repo repositories.AuthRepository, jwtCfg *config.JWTConfig) AuthUsecase {
	return &authUsecase{repo: repo, jwtCfg: jwtCfg}
}

func (u *authUsecase) GetAuth(req dto.AuthRequest) (*entities.Auth, error) {

	auth, err := u.repo.GetAuth(req)
	if err != nil {
		return nil, err
	}

	var result, accessToken, refreshToken, configKey string

	// 인증정보 없음.
	if auth.Id == "" {
		result = "fail"
	} else {
		result = "success"
		accessToken = getAccessToken(auth.Id, u.jwtCfg.AccessExp)
		refreshToken = getRefreshToken(auth.Id, u.jwtCfg.RefressExp)
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

func getAccessToken(id string, accessExp int) string {

	return ""
}

func getRefreshToken(id string, refreshExp int) string {

	return ""
}
