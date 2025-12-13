package repository

import (
	"batch/internal/domain/orgInfo/entity"
	"context"
)

type OrgInfoRepository interface {
	GetOrgInfo(ctx context.Context, org string) ([]entity.OrgInfoEntity, error)
	PutOrgInfoJson(ctx context.Context, org string, fileName string, orgJson string) error
}
