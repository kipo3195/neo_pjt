package repositories

import (
	"context"
	"fmt"
	"log"
	"org/entities"
	"org/models"

	"gorm.io/gorm"
)

const (
	DOMAIN = "domain"
	CODE   = "code"
)

type orgRepository struct {
	db *gorm.DB
}

type OrgRepository interface {
	CheckOrgHash(ctx context.Context, org string, hash string) (bool, bool, error)
	SaveDept(ctx context.Context, entity entities.CreateDeptEntity) (interface{}, error)
	DeleteDept(ctx context.Context, entity entities.DeleteDeptEntity) (interface{}, error)
	GetOrg(ctx context.Context, entity entities.GetOrgEntity) (*[]models.WorksOrg, error)
	SaveDeptUser(ctx context.Context, entity entities.CreateDeptUserEntity) (interface{}, error)
	DeleteDeptUser(txt context.Context, entity entities.DeleteDeptUserEntity) (interface{}, error)

	GetOrgLatestVersion(ctx context.Context, org string) (string, error)
	GetOrgDiffEvent(ctx context.Context, orgCode string, orgHash string) (models.OrgEvent, error)
}

func NewOrgRepository(db *gorm.DB) OrgRepository {
	return &orgRepository{db: db}
}

// 추가
func (r *orgRepository) SaveDept(ctx context.Context, entity entities.CreateDeptEntity) (interface{}, error) {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	models := r.ToWorksDeptsModel(entity)
	if err := tx.Create(&models).Error; err != nil {
		log.Println("[SaveDepartment] - DB error")
		tx.Rollback()
		return false, err
	}

	multiLangModels := r.ToWorksDeptsMultiLangModel(entity)
	if err := tx.Create(&multiLangModels).Error; err != nil {
		log.Println("[SaveDepartment Multi lang] - DB error")
		tx.Rollback()
		return false, err
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("[SaveDepartment] - Commit failed")
		return false, err
	}
	// DB 저장 성공
	fmt.Println("[SaveDepartment] success !")
	return true, nil
}

func (r *orgRepository) ToWorksDeptsModel(e entities.CreateDeptEntity) models.WorksDept {
	return models.WorksDept{
		DeptCode:        e.DeptCode,
		DeptOrg:         e.DeptOrg,
		ParentsDeptCode: e.ParentDeptCode,
	}
}

func (r *orgRepository) ToWorksDeptsMultiLangModel(e entities.CreateDeptEntity) models.WorksDeptMultiLang {
	return models.WorksDeptMultiLang{
		DeptCode: e.DeptCode,
		DeptOrg:  e.DeptOrg,
		KrLang:   e.KrLang,
		EnLang:   e.EnLang,
		JpLang:   e.JpLang,
		CnLang:   e.CnLang,
	}
}

// 삭제
func (r *orgRepository) DeleteDept(ctx context.Context, entity entities.DeleteDeptEntity) (interface{}, error) {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	fmt.Printf("부서 코드 : %s, 부서 org : %s \n", entity.DeptCode, entity.DeptOrg)

	// 첫 번째 삭제
	if err := tx.Where("dept_code = ? and dept_org = ?", entity.DeptCode, entity.DeptOrg).Delete(&models.WorksDept{}).Error; err != nil {
		log.Println("부서 메타데이터 삭제 실패:", err)
		tx.Rollback()
		return false, err
	}

	// 두 번째 삭제
	if err := tx.Where("dept_code = ? and dept_org = ?", entity.DeptCode, entity.DeptOrg).Delete(&models.WorksDeptMultiLang{}).Error; err != nil {
		log.Println("부서 멀티 랭기지 삭제 실패:", err)
		tx.Rollback()
		return false, err
	}

	// 트랜잭션 반영
	tx.Commit()

	// DB 저장 성공
	fmt.Println("[DeleteDepartment] success !")
	return true, nil
}

func (r *orgRepository) GetOrg(ctx context.Context, entity entities.GetOrgEntity) (*[]models.WorksOrg, error) {

	var orgTree *[]models.WorksOrg
	viewSql := `SELECT * FROM org.vw_dept_and_user_tree where org = ?`
	err := r.db.Raw(viewSql, entity.OrgCode).Scan(&orgTree).Error

	if err != nil {
		log.Println("[GetOrg] - No record found or DB error")
		return nil, err
	}

	return orgTree, nil
}

func (r *orgRepository) SaveDeptUser(ctx context.Context, entity entities.CreateDeptUserEntity) (interface{}, error) {
	// 트랜잭션 시작
	// tx := r.db.WithContext(ctx).Begin()
	// if tx.Error != nil {
	// 	return false, tx.Error
	// }

	models := toWorksDeptUser(entity)
	if err := r.db.Create(&models).Error; err != nil {
		log.Println("[SaveDeptUser] - DB error")
		// tx.Rollback()
		return false, err
	}

	// if err := tx.Commit().Error; err != nil {
	// 	log.Println("[SaveDeptUser] - Commit failed")
	// 	return false, err
	// }
	// DB 저장 성공
	fmt.Println("[SaveDeptUser] success !")
	return true, nil
}

func toWorksDeptUser(entity entities.CreateDeptUserEntity) *models.WorksDeptUser {
	return &models.WorksDeptUser{
		DeptCode:             entity.DeptCode,
		DeptOrg:              entity.DeptOrg,
		UserHash:             entity.UserHash,
		PositionCode:         entity.PositionCode,
		RoleCode:             entity.RoleCode,
		IsConcurrentPosition: entity.IsConcurrentPosition,
		UpdateHash:           entity.UpdateHash,
	}
}

func (r *orgRepository) DeleteDeptUser(txt context.Context, entity entities.DeleteDeptUserEntity) (interface{}, error) {

	result := r.db.Model(&models.WorksDeptUser{}).Where("dept_org = ? AND dept_code = ? AND user_hash = ? ", entity.DeptOrg, entity.DeptCode, entity.UserHash).Update("use_yn", "N")
	return nil, result.Error
}

func (r *orgRepository) CheckOrgHash(ctx context.Context, org string, hash string) (bool, bool, error) {

	var count int64
	var events models.OrgEvent

	// hash 검증 'req'데이터를 '_' 기준으로 split하면 [0]은 org code, [1]은 hash

	result := r.db.Model(events).
		Where("update_hash >= ? AND update_hash <= ? AND org_code = ?", hash, 9999999999999999, org).
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
	err := r.db.Where("org_code", org).Order("hash DESC").First(&result).Error
	if err != nil {
		log.Println("Error fetching latest hash:", err)
	}
	fmt.Printf("org : %s Latest hash record : %s \n", org, result.UpdateHash)
	return result.UpdateHash, nil

}

func (r *orgRepository) GetOrgDiffEvent(ctx context.Context, orgCode string, orgHash string) (models.OrgEvent, error) {

	var events models.OrgEvent

	err := r.db.Where("update_hash >= ? AND update_hash <= ? AND org_code = ?", orgHash, 9999999999999999, orgCode).Find(&events).Error

	if err != nil {
		return events, err
	}

	return events, nil
}
