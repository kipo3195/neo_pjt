package repositories

import (
	"context"
	"fmt"
	"log"
	"org/consts"
	"org/models"
	"time"

	"gorm.io/gorm"
)

func generateUpdateHash() string {
	return time.Now().Format(consts.YYYYMMDDHHMSS)
}

type orgRepository struct {
	db *gorm.DB
}

type OrgRepository interface {
	CheckOrgHash(ctx context.Context, org string, hash string) (bool, bool, error)

	GetOrgLatestVersion(ctx context.Context, org string) (string, error)
	GetOrgDiffEvent(ctx context.Context, orgCode string, orgHash string) ([]models.OrgEvent, error)
}

func NewOrgRepository(db *gorm.DB) OrgRepository {
	return &orgRepository{db: db}
}

func (r *orgRepository) CheckOrgHash(ctx context.Context, org string, hash string) (bool, bool, error) {

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

func (r *orgRepository) GetOrgLatestVersion(ctx context.Context, org string) (string, error) {

	var result models.OrgEventHash
	err := r.db.Where("org_code = ?", org).Order("update_hash DESC").First(&result).Error
	if err != nil {
		// 결과가 없을때
		log.Println("Error fetching latest hash:", err)
		return "", err
	}
	fmt.Printf("org : %s Latest hash record : %s \n", org, result.UpdateHash)
	return result.UpdateHash, nil

}

func (r *orgRepository) GetOrgDiffEvent(ctx context.Context, orgCode string, orgHash string) ([]models.OrgEvent, error) {

	var events []models.OrgEvent

	err := r.db.Where("org_code = ? AND update_hash > ?", orgCode, orgHash).Find(&events).Error

	if err != nil {
		return nil, err
	}

	return events, nil
}
