package server

import (
	"context"
	"log"
	"org/models"

	"gorm.io/gorm"
)

type orgRepository struct {
	db *gorm.DB
}

type OrgRepository interface {
	PutOrgEventHash(ctx context.Context, org string, hash string) (bool, error)
	GetOrg(ctx context.Context, orgCode string) ([]models.WorksOrg, error)
}

func NewOrgRepository(db *gorm.DB) OrgRepository {
	return &orgRepository{
		db: db,
	}
}

func (r *orgRepository) PutOrgEventHash(ctx context.Context, org string, hash string) (bool, error) {

	models := toOrgEventHashModel(org, hash)
	if err := r.db.Create(&models).Error; err != nil {
		log.Println("[PutOrgEventHash] - DB error")
		return false, err
	}
	return true, nil
}

func toOrgEventHashModel(org string, hash string) models.OrgEventHash {
	return models.OrgEventHash{
		OrgCode:    org,
		UpdateHash: hash,
	}
}

func (r *orgRepository) GetOrg(ctx context.Context, orgCode string) ([]models.WorksOrg, error) {

	var orgTree []models.WorksOrg
	viewSql := `SELECT * FROM org.vw_dept_and_user_tree where org = ?`
	err := r.db.Raw(viewSql, orgCode).Scan(&orgTree).Error

	if err != nil {
		log.Println("[GetOrg] - No record found or DB error")
		return nil, err
	}

	return orgTree, nil
}
