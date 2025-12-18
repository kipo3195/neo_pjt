package repository

import "batch/internal/domain/userDetail/repository"

type userDetailApiRepositoryImpl struct {
	domain string
}

func NewUserDetailApiRepository(domain string) repository.UserDetailApiRepository {
	return &userDetailApiRepositoryImpl{
		domain: domain,
	}
}
