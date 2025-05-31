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
	SaveDepartment(ctx context.Context, entity entities.CreateDepartmentEntity) (interface{}, error)
	DeleteDepartment(ctx context.Context, entity entities.DeleteDepartmentEntity) (interface{}, error)
}

func NewOrgRepository(db *gorm.DB) OrgRepository {
	return &orgRepository{db: db}
}

// 추가
func (r *orgRepository) SaveDepartment(ctx context.Context, entity entities.CreateDepartmentEntity) (interface{}, error) {

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

func (r *orgRepository) ToWorksDeptsModel(e entities.CreateDepartmentEntity) models.WorksDept {
	return models.WorksDept{
		DeptCode:        e.DeptCode,
		DeptOrg:         e.DeptOrg,
		ParentsDeptCode: e.ParentDeptCode,
	}
}

func (r *orgRepository) ToWorksDeptsMultiLangModel(e entities.CreateDepartmentEntity) models.WorksDeptMultiLang {
	return models.WorksDeptMultiLang{
		DeptCode:   e.DeptCode,
		DeptOrg:    e.DeptOrg,
		DeptNameKr: e.DeptNameKr,
		DeptNameEn: e.DeptNameEn,
		DeptNameCn: e.DeptNameCn,
		DeptNameJp: e.DeptNameJp,
	}
}

// 삭제
func (r *orgRepository) DeleteDepartment(ctx context.Context, entity entities.DeleteDepartmentEntity) (interface{}, error) {

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
