package repository

import (
	"context"
	"fmt"
	"log"
	"org/internal/domain/department/entity"
	"org/internal/domain/department/repository"
	"org/internal/infrastructure/model"
	utils "org/internal/infrastructure/util"
	"time"

	"gorm.io/gorm"
)

type departmentRepositoryImpl struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) repository.DepartmentRepository {
	return &departmentRepositoryImpl{
		db: db,
	}
}

func DepartmentMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.WorksDept{})
	db.AutoMigrate(&model.WorksDeptMultiLang{})
	db.AutoMigrate(&model.PositionMultiLang{})
	db.AutoMigrate(&model.RoleMultiLang{})
	db.AutoMigrate(&model.WorksDeptUser{})
	db.AutoMigrate(&model.WorksUserMultiLang{})
}

// 추가
func (r *departmentRepositoryImpl) PutDept(ctx context.Context, entity entity.CreateDeptEntity) (interface{}, error) {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	// 부서 추가
	worksDept := toWorksDeptsModel(entity)
	if err := tx.Create(&worksDept).Error; err != nil {
		log.Println("[PutDept] - DB error")
		tx.Rollback()
		return false, err
	}

	// 부서의 다국어 추가
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
	log.Println("[PutDept] success !")
	return true, nil
}

func toWorksDeptsModel(e entity.CreateDeptEntity) model.WorksDept {
	return model.WorksDept{
		DeptCode:        e.DeptCode,
		DeptOrg:         e.DeptOrg,
		ParentsDeptCode: e.ParentDeptCode,
		Header:          e.Header,
	}
}

func toWorksDeptsMultiLangModel(e entity.CreateDeptEntity) model.WorksDeptMultiLang {
	return model.WorksDeptMultiLang{
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

func toOrgCreateEventModel(e entity.CreateDeptEntity) model.OrgEvent {
	return model.OrgEvent{
		EventType:  "C",
		Id:         e.DeptCode,
		Kind:       "0",
		OrgCode:    e.DeptOrg,
		UpdateHash: utils.GenerateUpdateHash(),
	}
}

// 삭제
func (r *departmentRepositoryImpl) DeleteDept(ctx context.Context, entity entity.DeleteDeptEntity) (interface{}, error) {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	fmt.Printf("부서 코드 : %s, 부서 org : %s \n", entity.DeptCode, entity.DeptOrg)

	// 첫 번째 삭제
	if err := tx.Where("dept_code = ? and dept_org = ?", entity.DeptCode, entity.DeptOrg).Delete(&model.WorksDept{}).Error; err != nil {
		log.Println("부서 메타데이터 삭제 실패:", err)
		tx.Rollback()
		return false, err
	}

	// 두 번째 삭제
	if err := tx.Where("dept_code = ? and dept_org = ?", entity.DeptCode, entity.DeptOrg).Delete(&model.WorksDeptMultiLang{}).Error; err != nil {
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
	log.Println("[DeleteDepartment] success !")
	return true, nil
}

func toOrgDeleteEventModel(e entity.DeleteDeptEntity) model.OrgEvent {
	return model.OrgEvent{
		EventType:  "D",
		Id:         e.DeptCode,
		Kind:       "0",
		OrgCode:    e.DeptOrg,
		UpdateHash: utils.GenerateUpdateHash(),
	}
}

func (r *departmentRepositoryImpl) PutDeptUser(ctx context.Context, entity entity.CreateDeptUserEntity) (interface{}, error) {
	// 트랜잭션 시작
	// tx := r.db.WithContext(ctx).Begin()
	// if tx.Error != nil {
	// 	return false, tx.Error
	// }

	model := toWorksDeptUser(entity)
	if err := r.db.Create(&model).Error; err != nil {
		log.Println("[PutDeptUser] - DB error")
		// tx.Rollback()
		return false, err
	}

	// if err := tx.Commit().Error; err != nil {
	// 	log.Println("[SaveDeptUser] - Commit failed")
	// 	return false, err
	// }
	// DB 저장 성공
	log.Println("[PutDeptUser] success !")
	return true, nil
}

func toWorksDeptUser(entity entity.CreateDeptUserEntity) *model.WorksDeptUser {
	return &model.WorksDeptUser{
		DeptCode:             entity.DeptCode,
		DeptOrg:              entity.DeptOrg,
		UserHash:             entity.UserHash,
		PositionCode:         entity.PositionCode,
		RoleCode:             entity.RoleCode,
		IsConcurrentPosition: entity.IsConcurrentPosition,
		UpdateHash:           entity.UpdateHash,
	}
}

func (r *departmentRepositoryImpl) DeleteDeptUser(txt context.Context, entity entity.DeleteDeptUserEntity) (interface{}, error) {

	result := r.db.Model(&model.WorksDeptUser{}).Where("dept_org = ? AND dept_code = ? AND user_hash = ? ", entity.DeptOrg, entity.DeptCode, entity.UserHash).Update("use_yn", "N")
	return nil, result.Error
}

func (r *departmentRepositoryImpl) CreateDeptTree(ctx context.Context, e entity.WorksDeptEntity) error {
	deptModel := model.WorksDept{
		DeptCode:        e.DeptCode,
		DeptOrg:         e.DeptOrg,
		ParentsDeptCode: e.ParentsDeptCode,
	}

	// DeptCreateDate가 비어있다면 현재 시각으로 설정
	if deptModel.DeptCreateDate.IsZero() {
		deptModel.DeptCreateDate = time.Now()
	}

	if err := r.db.WithContext(ctx).Create(&deptModel).Error; err != nil {
		return err
	}

	return nil
}

func (r *departmentRepositoryImpl) GetDepts(ctx context.Context, org string) ([]entity.WorksDeptEntity, error) {

	var worksOrgs []model.WorksDept

	err := r.db.Raw(`SELECT * FROM org.works_dept where dept_org = ?`, org).Scan(&worksOrgs).Error

	if err != nil {
		log.Println("[GetDepts] - No record found or DB error")
		return nil, err
	}

	temp := make([]entity.WorksDeptEntity, len(worksOrgs))

	for i := 0; i < len(worksOrgs); i++ {
		log.Println(worksOrgs[i].DeptCode, " : ", worksOrgs[i].DeptOrg)
		temp[i] = entity.WorksDeptEntity{
			DeptCode: worksOrgs[i].DeptCode,
			DeptOrg:  worksOrgs[i].DeptOrg,
		}
	}

	return temp, nil
}

func (r *departmentRepositoryImpl) PutWorksDeptMultiLang(ctx context.Context, e entity.CreateMultiLangEntity) error {

	log.Println("insert > ", e.DeptCode, " : ", e.DeptOrg)

	deptModel := model.WorksDeptMultiLang{
		DeptCode: e.DeptCode,
		DeptOrg:  e.DeptOrg,
		KoLang:   e.KoLang,
		EnLang:   e.EnLang,
		ZhLang:   e.ZhLang,
		JpLang:   e.JpLang,
		RuLang:   e.RuLang,
		ViLang:   e.ViLang,
		DefLang:  e.KoLang,
	}

	if err := r.db.WithContext(ctx).Create(&deptModel).Error; err != nil {
		return err
	}

	return nil
}
