package repository

import "admin/internal/domain/userAuthRegister/repository"

type userAuthRegisterRepositoryImpl struct {
}

func NewUserAuthRegisterRepository() repository.UserAuthRegisterRepository {
	return &userAuthRegisterRepositoryImpl{}
}
