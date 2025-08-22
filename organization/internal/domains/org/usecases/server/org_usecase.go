package server

import (
	repositories "org/internal/domains/org/repositories/server"
)

type orgUsecase struct {
	repository repositories.OrgRepository
}

type OrgUsecase interface {
}

func NewOrgUsecase(repository repositories.OrgRepository) OrgUsecase {
	return &orgUsecase{
		repository: repository,
	}
}
