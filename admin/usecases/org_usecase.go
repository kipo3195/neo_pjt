package usecases

import (
	"admin/repositories"
)

type orgUsecase struct {
	repo repositories.OrgRepository
}

type OrgUsecase interface {
}

func NewOrgUsecase(repo repositories.OrgRepository) OrgUsecase {
	return &orgUsecase{repo: repo}
}
