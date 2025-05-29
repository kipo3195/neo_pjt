package usecases

import (
	"admin/repositories"
)

type adminUsecase struct {
	repo repositories.AdminRepository
}

type AdminUsecase interface {
}

func NewAdminUsecase(repo repositories.AdminRepository) AdminUsecase {
	return &adminUsecase{repo: repo}
}
