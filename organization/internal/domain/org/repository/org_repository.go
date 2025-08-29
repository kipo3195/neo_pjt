package repository

import (
	"context"
	"org/models"
)

type OrgRepository interface {
	CheckOrgHash(ctx context.Context, org string, hash string) (bool, bool, error)
	GetOrgLatestVersion(ctx context.Context, orgCode string) (string, error)
	GetOrgDiffEvent(ctx context.Context, orgCode string, orgHash string) ([]models.OrgEvent, error)
	PutOrgEventHash(ctx context.Context, org string, hash string) (bool, error)
	GetOrg(ctx context.Context, orgCode string) ([]models.WorksOrg, error)
}
