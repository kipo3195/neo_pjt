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
	GetOrg(ctx context.Context, entity entities.GetOrgEntity) (*[]models.WorksOrg, error)
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
		DeptCode: e.DeptCode,
		DeptOrg:  e.DeptOrg,
		KrLang:   e.KrLang,
		EnLang:   e.EnLang,
		JpLang:   e.JpLang,
		CnLang:   e.CnLang,
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

func (r *orgRepository) GetOrg(ctx context.Context, entity entities.GetOrgEntity) (*[]models.WorksOrg, error) {

	var orgTree *[]models.WorksOrg
	treeSql := `WITH RECURSIVE dept_tree AS (
			SELECT 
				dept_code,
				parent_dept_code,
				dept_update_hash
			FROM works_dept
			WHERE parent_dept_code = 'root'
			UNION ALL
			SELECT 
				d.dept_code,
				d.parent_dept_code,
				d.dept_update_hash
			FROM works_dept d
			INNER JOIN dept_tree dt ON d.parent_dept_code = dt.dept_code
		) SELECT a.dept_code, a.parent_dept_code, b.kr_lang, b.en_lang, b.cn_lang, b.jp_lang, a.dept_update_hash 
		FROM dept_tree as a join works_dept_multi_lang as b on a.dept_code = b.dept_code ;`

	err := r.db.Raw(treeSql).Scan(&orgTree).Error

	if err != nil {
		log.Println("[GetOrg] - No record found or DB error")
		return nil, err
	}

	return orgTree, nil
}
