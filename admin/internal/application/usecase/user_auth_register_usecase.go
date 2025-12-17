package usecase

import (
	"admin/internal/application/usecase/input"
	"admin/internal/consts"
	"admin/internal/domain/userAuthRegister/entity"
	"admin/internal/domain/userAuthRegister/repository"
	"context"
	"log"
)

type userAuthRegisterUsecase struct {
	repo    repository.UserAuthRegisterRepository
	apiRepo repository.UserAuthRegisterApiRepository
}

type UserAuthRegisterUsecase interface {
	UserAuthRegisterInAuth(ctx context.Context, in input.UserAuthRegisterInput) error
}

func NewUserAuthRegisterUsecase(repo repository.UserAuthRegisterRepository, apiRepo repository.UserAuthRegisterApiRepository) UserAuthRegisterUsecase {

	return &userAuthRegisterUsecase{
		repo:    repo,
		apiRepo: apiRepo,
	}
}

func (r *userAuthRegisterUsecase) UserAuthRegisterInAuth(ctx context.Context, in input.UserAuthRegisterInput) error {

	en := make([]entity.UserAuthRegisterEntity, 0)

	for _, i := range in.ServiceUser {

		temp := entity.UserAuthRegisterEntity{
			UserId:   i.UserId,
			UserHash: i.UserHash,
			UserAuth: i.UserAuth,
			Salt:     i.Salt,
		}

		en = append(en, temp)
	}

	if len(en) > 0 {
		result, err := r.apiRepo.UserAuthRegisterInAuth(ctx, en)
		if err != nil {
			log.Println("[UserAuthRegisterInAuth] err : ", err)
			return err
		} else if result == "fail" {
			log.Println("[UserAuthRegisterInAuth] result fail.")
			return consts.ErrUserAuthRegisterFail
		}
		return nil
	} else {
		log.Println("[UserAuthRegisterInAuth] entity size = 0")
		return consts.ErrUserAuthRegisterEntitySizeError
	}
}
