package client

import "org/repositories"

type userUsecase struct {
	repository respositories.UserRepository
}

type UserUsecase interface {
}

func NewUserUsecase(repository repositories.UserRepository)
