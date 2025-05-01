package usecases

import "message/repositories"

type authUsecase struct {
	repo repositories.AuthRepository
}
type AuthUsecase interface {
	AuthenticateToken(token string) (bool, error)
}

func NewAuthUsecase(repo repositories.AuthRepository) AuthUsecase {
	return &authUsecase{repo: repo}
}

func (uc *authUsecase) AuthenticateToken(token string) (bool, error) {
	// 토큰 검증 로직

	if token == "" {
		return false, nil
	}

	return true, nil
}
