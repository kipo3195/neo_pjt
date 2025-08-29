package repository

import (
	"context"
	"fmt"
	"log"
	"org/models"

	"org/internal/domain/org/repository"

	"gorm.io/gorm"
)

type orgRepositoryImpl struct {
	db *gorm.DB
}

func NewOrgRepository(db *gorm.DB) repository.OrgRepository {
	return &orgRepositoryImpl{
		db: db,
	}
}

func (r *orgRepositoryImpl) CheckOrgHash(ctx context.Context, org string, hash string) (bool, bool, error) {

	var count int64
	var events models.OrgEvent

	// hash 검증 'req'데이터를 '_' 기준으로 split하면 [0]은 org code, [1]은 hash

	result := r.db.Model(events).
		Where("update_hash > ? AND org_code = ?", hash, org).
		Count(&count)

	if result.Error != nil {
		// 에러
		log.Printf("Count query failed: %v", result.Error)
		return false, false, result.Error
	}

	// 50개 이상이면 파일로 처리 필요.
	if count >= 50 {
		return true, false, nil
	} else if count == 0 {
		return false, false, nil
	} else {
		return false, true, nil
	}
}

func (r *orgRepositoryImpl) GetOrgLatestVersion(ctx context.Context, orgCode string) (string, error) {

	var result models.OrgEventHash
	err := r.db.Where("org_code = ?", orgCode).Order("update_hash DESC").First(&result).Error
	if err != nil {
		// 결과가 없을때
		log.Println("Error fetching latest hash:", err)
		return "", err
	}
	fmt.Printf("org : %s Latest hash record : %s \n", orgCode, result.UpdateHash)
	return result.UpdateHash, nil

}

func (r *orgRepositoryImpl) GetOrgDiffEvent(ctx context.Context, orgCode string, orgHash string) ([]models.OrgEvent, error) {

	var events []models.OrgEvent

	err := r.db.Where("org_code = ? AND update_hash > ?", orgCode, orgHash).Find(&events).Error

	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *orgRepositoryImpl) PutOrgEventHash(ctx context.Context, org string, hash string) (bool, error) {

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

func (r *orgRepositoryImpl) GetOrg(ctx context.Context, orgCode string) ([]models.WorksOrg, error) {

	var orgTree []models.WorksOrg
	viewSql := `SELECT * FROM org.vw_dept_and_user_tree where org = ?`
	err := r.db.Raw(viewSql, orgCode).Scan(&orgTree).Error

	if err != nil {
		log.Println("[GetOrg] - No record found or DB error")
		return nil, err
	}

	return orgTree, nil
}
