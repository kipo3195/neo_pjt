package client

import (
	repositories "org/internal/domains/org/repositories/client"
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
