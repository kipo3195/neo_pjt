package repositories

import (
	"context"
	"fmt"
	"log"
	"org/entities"
	"org/models"

	"gorm.io/gorm"
)

type serverRepository struct {
	db *gorm.DB
}

type ServerRepository interface {
	PutDept(ctx context.Context, entity entities.CreateDeptEntity) (interface{}, error)
	DeleteDept(ctx context.Context, entity entities.DeleteDeptEntity) (interface{}, error)

	PutDeptUser(ctx context.Context, entity entities.CreateDeptUserEntity) (interface{}, error)
	DeleteDeptUser(ctx context.Context, entity entities.DeleteDeptUserEntity) (interface{}, error)

	PutOrgEventHash(ctx context.Context, org string, hash string) (bool, error)

	GetOrg(ctx context.Context, entity entities.GetOrgEntity) ([]models.WorksOrg, error)
}

func NewServerRepository(db *gorm.DB) ServerRepository {
	return &serverRepository{db: db}
}

// 추가
func (r *serverRepository) PutDept(ctx context.Context, entity entities.CreateDeptEntity) (interface{}, error) {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	worksDept := toWorksDeptsModel(entity)
	if err := tx.Create(&worksDept).Error; err != nil {
		log.Println("[PutDept] - DB error")
		tx.Rollback()
		return false, err
	}

	worksDeptMultiLang := toWorksDeptsMultiLangModel(entity)
	if err := tx.Create(&worksDeptMultiLang).Error; err != nil {
		log.Println("[PutDept Multi lang] - DB error")
		tx.Rollback()
		return false, err
	}

	// org_event에 추가.
	orgEventModel := toOrgCreateEventModel(entity)
	if err := tx.Create(&orgEventModel).Error; err != nil {
		log.Println("[PutDept org event] - DB error")
		tx.Rollback()
		return false, err
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("[PutDept] - Commit failed")
		return false, err
	}
	// DB 저장 성공
	fmt.Println("[PutDept] success !")
	return true, nil
}

func toWorksDeptsModel(e entities.CreateDeptEntity) models.WorksDept {
	return models.WorksDept{
		DeptCode:        e.DeptCode,
		DeptOrg:         e.DeptOrg,
		ParentsDeptCode: e.ParentDeptCode,
	}
}

func toWorksDeptsMultiLangModel(e entities.CreateDeptEntity) models.WorksDeptMultiLang {
	return models.WorksDeptMultiLang{
		DeptCode: e.DeptCode,
		DeptOrg:  e.DeptOrg,
		KoLang:   e.KoLang,
		EnLang:   e.EnLang,
		JpLang:   e.JpLang,
		ZhLang:   e.ZhLang,
		RuLang:   e.RuLang,
		ViLang:   e.ViLang,
	}
}

func toOrgCreateEventModel(e entities.CreateDeptEntity) models.OrgEvent {
	return models.OrgEvent{
		EventType:  "C",
		Id:         e.DeptCode,
		Kind:       "0",
		OrgCode:    e.DeptOrg,
		UpdateHash: generateUpdateHash(),
	}
}

// 삭제
func (r *serverRepository) DeleteDept(ctx context.Context, entity entities.DeleteDeptEntity) (interface{}, error) {

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

	// org_event에 추가.
	orgEventModel := toOrgDeleteEventModel(entity)
	if err := tx.Create(&orgEventModel).Error; err != nil {
		log.Println("[SaveDepartment org event] - DB error")
		tx.Rollback()
		return false, err
	}

	// 트랜잭션 반영
	tx.Commit()

	// DB 저장 성공
	fmt.Println("[DeleteDepartment] success !")
	return true, nil
}

func toOrgDeleteEventModel(e entities.DeleteDeptEntity) models.OrgEvent {
	return models.OrgEvent{
		EventType:  "D",
		Id:         e.DeptCode,
		Kind:       "0",
		OrgCode:    e.DeptOrg,
		UpdateHash: generateUpdateHash(),
	}
}

func (r *serverRepository) PutDeptUser(ctx context.Context, entity entities.CreateDeptUserEntity) (interface{}, error) {
	// 트랜잭션 시작
	// tx := r.db.WithContext(ctx).Begin()
	// if tx.Error != nil {
	// 	return false, tx.Error
	// }

	models := toWorksDeptUser(entity)
	if err := r.db.Create(&models).Error; err != nil {
		log.Println("[PutDeptUser] - DB error")
		// tx.Rollback()
		return false, err
	}

	// if err := tx.Commit().Error; err != nil {
	// 	log.Println("[SaveDeptUser] - Commit failed")
	// 	return false, err
	// }
	// DB 저장 성공
	fmt.Println("[PutDeptUser] success !")
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

func (r *serverRepository) DeleteDeptUser(txt context.Context, entity entities.DeleteDeptUserEntity) (interface{}, error) {

	result := r.db.Model(&models.WorksDeptUser{}).Where("dept_org = ? AND dept_code = ? AND user_hash = ? ", entity.DeptOrg, entity.DeptCode, entity.UserHash).Update("use_yn", "N")
	return nil, result.Error
}

func (r *serverRepository) PutOrgEventHash(ctx context.Context, org string, hash string) (bool, error) {

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

func (r *serverRepository) GetOrg(ctx context.Context, entity entities.GetOrgEntity) ([]models.WorksOrg, error) {

	var orgTree []models.WorksOrg
	viewSql := `SELECT * FROM org.vw_dept_and_user_tree where org = ?`
	err := r.db.Raw(viewSql, entity.OrgCode).Scan(&orgTree).Error

	if err != nil {
		log.Println("[GetOrg] - No record found or DB error")
		return nil, err
	}

	return orgTree, nil
}
