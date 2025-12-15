package usecase

import (
	"common/internal/domain/org/repository"
	"common/internal/infrastructure/storage"
)

type orgUsecase struct {
	repo       repository.OrgRepository
	orgStorage storage.OrgStorage
}

type OrgUsecase interface {
	GetWorksOrgCode() []string
}

func NewOrgUsecase(repo repository.OrgRepository, orgStorage storage.OrgStorage) OrgUsecase {
	return &orgUsecase{
		repo:       repo,
		orgStorage: orgStorage,
	}
}

func (r *orgUsecase) GetWorksOrgCode() []string {

	return r.orgStorage.GetOrgInfo()
}
