package usecase

import (
	"context"
	"fmt"
	"test/internal/application/usecase/input"
	"test/internal/infrastructure/config"
)

type loginUsecase struct {
}
type LoginUsecase interface {
	PutLogin(ctx context.Context, input input.PutLoginInput)
}

func NewLoginUsecase(sfg *config.ServerConfig) LoginUsecase {

	return &loginUsecase{}
}

func (r *loginUsecase) PutLogin(ctx context.Context, input input.PutLoginInput) {

	fmt.Println("로깅 테스트")

}
