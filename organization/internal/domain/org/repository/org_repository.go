package repository

import (
	"context"
	"org/internal/domain/org/entity"
)

type OrgRepository interface {
	CheckOrgHash(ctx context.Context, org string, hash string) (bool, bool, error)
	GetOrgLatestVersion(ctx context.Context, orgCode string) (string, error)
	GetOrgDiffEvent(ctx context.Context, orgCode string, orgHash string) ([]entity.OrgEventEntity, error)
	PutOrgEventHash(ctx context.Context, org string, hash string) (bool, error)
	GetOrg(ctx context.Context, orgCode string) ([]entity.WorksOrg, error)
	RegistOrgBatch(ctx context.Context, dept []entity.WorksOrg, user []entity.WorksOrg) error
}
