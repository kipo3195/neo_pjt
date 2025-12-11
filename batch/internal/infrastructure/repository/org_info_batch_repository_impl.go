package repository

import "batch/internal/domain/orgInfo/repository"

type orgInfoBatchRepository struct {
}

func NewOrgInfoBatchRepository() repository.OrgInfoBatchRepository {

	return &orgInfoBatchRepository{}

}
