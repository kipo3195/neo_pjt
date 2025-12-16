package repository

import (
	"context"
	"fmt"
	"log"
	"org/internal/domain/org/entity"
	"org/internal/domain/org/repository"
	"org/internal/infrastructure/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type orgRepositoryImpl struct {
	db *gorm.DB
}

func NewOrgRepository(db *gorm.DB) repository.OrgRepository {
	return &orgRepositoryImpl{
		db: db,
	}
}

func OrgMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.OrgEvent{})
	db.AutoMigrate(&model.OrgEventHash{})
	db.AutoMigrate(&model.WorksDeptUser{})
	db.AutoMigrate(&model.WorksOrgCode{})
}

func (r *orgRepositoryImpl) CheckOrgHash(ctx context.Context, org string, hash string) (bool, bool, error) {

	var count int64
	var events model.OrgEvent

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

	var result model.OrgEventHash
	err := r.db.Where("org_code = ?", orgCode).Order("update_hash DESC").First(&result).Error
	if err != nil {
		// 결과가 없을때
		log.Println("Error fetching latest hash:", err)
		return "", err
	}
	fmt.Printf("org : %s Latest hash record : %s \n", orgCode, result.UpdateHash)
	return result.UpdateHash, nil

}

func (r *orgRepositoryImpl) GetOrgDiffEvent(ctx context.Context, orgCode string, orgHash string) ([]entity.OrgEventEntity, error) {

	var events []model.OrgEvent

	err := r.db.Where("org_code = ? AND update_hash > ?", orgCode, orgHash).Find(&events).Error

	if err != nil {
		return nil, err
	}

	return convertToOrgEventEntities(events), nil
}

// 변환 함수
func convertToOrgEventEntities(events []model.OrgEvent) []entity.OrgEventEntity {
	entities := make([]entity.OrgEventEntity, 0, len(events))
	for _, e := range events {
		entity := entity.OrgEventEntity{
			Id:         e.Id,
			EventType:  e.EventType,
			Kind:       e.Kind,
			OrgCode:    e.OrgCode,
			UpdateHash: e.UpdateHash,
		}
		entities = append(entities, entity)
	}
	return entities
}

func (r *orgRepositoryImpl) PutOrgEventHash(ctx context.Context, org string, hash string) (bool, error) {

	model := toOrgEventHashModel(org, hash)
	if err := r.db.Create(&model).Error; err != nil {
		log.Println("[PutOrgEventHash] - DB error")
		return false, err
	}
	return true, nil
}

func toOrgEventHashModel(org string, hash string) model.OrgEventHash {
	return model.OrgEventHash{
		OrgCode:    org,
		UpdateHash: hash,
	}
}

func (r *orgRepositoryImpl) GetOrg(ctx context.Context, orgCode string) ([]entity.WorksOrg, error) {

	var worksOrgs []model.WorksOrg
	viewSql := `SELECT * FROM org.vw_dept_and_user_tree where org = ?`
	err := r.db.Raw(viewSql, orgCode).Scan(&worksOrgs).Error

	if err != nil {
		log.Println("[GetOrg] - No record found or DB error")
		return nil, err
	}

	// 파싱 로직 추가 var orgTree []model.WorksOrg -> []entity.WorksOrg

	return convertWorksOrgToEntity(worksOrgs), nil
}

// 변환 함수
func convertWorksOrgToEntity(models []model.WorksOrg) []entity.WorksOrg {
	entities := make([]entity.WorksOrg, 0, len(models))
	for _, m := range models {
		entity := entity.WorksOrg{
			Org:            m.Org,
			DeptCode:       m.DeptCode,
			ParentDeptCode: m.ParentDeptCode,
			KoLang:         m.KoLang,
			EnLang:         m.EnLang,
			ZhLang:         m.ZhLang,
			JpLang:         m.JpLang,
			RuLang:         m.RuLang,
			ViLang:         m.ViLang,
			UpdateHash:     m.UpdateHash,
			Kind:           m.Kind,
			UserHash:       m.UserHash,
			UserId:         m.UserId,
			Header:         m.Header,
			Description:    m.Description,
		}
		entities = append(entities, entity)
	}
	return entities
}

func (r *orgRepositoryImpl) RegistOrgBatch(ctx context.Context, dept []entity.WorksOrg, user []entity.WorksOrg) error {

	tx := r.db.Begin()

	/* 부서 */
	worksDept := make([]model.WorksDept, 0)
	worksDeptMultiLang := make([]model.WorksDeptMultiLang, 0)

	for _, d := range dept {

		worksDept = append(worksDept, model.WorksDept{
			DeptCode:        d.DeptCode,
			DeptOrg:         d.Org,
			ParentsDeptCode: d.ParentDeptCode,
			UpdateHash:      d.UpdateHash,
		})

		worksDeptMultiLang = append(worksDeptMultiLang, model.WorksDeptMultiLang{
			DeptCode: d.DeptCode,
			DeptOrg:  d.Org,
			KoLang:   d.KoLang,
			EnLang:   d.EnLang,
			ZhLang:   d.ZhLang,
			JpLang:   d.JpLang,
			RuLang:   d.RuLang,
			ViLang:   d.ViLang,
			DefLang:  d.KoLang,
		})
	}

	if len(worksDept) > 0 {

		err := tx.Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "dept_org"},
					{Name: "dept_code"},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"update_hash",
				}),
			},
		).Create(&worksDept).Error

		if err != nil {
			return err
		}

		err = tx.Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "dept_org"},
					{Name: "dept_code"},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"ko_lang",
					"en_lang",
					"zh_lang",
					"jp_lang",
					"ru_lang",
					"vi_lang",
					"def_lang",
				}),
			},
		).Create(&worksDeptMultiLang).Error

		if err != nil {
			return err
		}
	}

	/* 사용자 */

	userDetail := make([]model.UserDetail, 0)
	worksDeptUser := make([]model.WorksDeptUser, 0)
	worksUserMultiLang := make([]model.WorksUserMultiLang, 0)

	for _, u := range user {

		userDetail = append(userDetail, model.UserDetail{
			UserHash: u.UserHash,
		})

		worksDeptUser = append(worksDeptUser, model.WorksDeptUser{
			DeptCode:   u.DeptCode,
			DeptOrg:    u.Org,
			UserHash:   u.UserHash,
			UpdateHash: u.UpdateHash,
		})

		worksUserMultiLang = append(worksUserMultiLang, model.WorksUserMultiLang{
			UserHash: u.UserHash,
			KoLang:   u.KoLang,
			EnLang:   u.EnLang,
			ZhLang:   u.ZhLang,
			JpLang:   u.JpLang,
			RuLang:   u.RuLang,
			ViLang:   u.ViLang,
			DefLang:  u.KoLang,
		})

	}

	if len(userDetail) > 0 {

		err := tx.Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "user_hash"},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"update_hash",
					"user_phone_num",
					"user_email",
				}),
			},
		).Create(&userDetail).Error

		if err != nil {
			return err
		}

		err = tx.Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "dept_code"},
					{Name: "user_hash"},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"update_hash",
					"role_code",
					"rank_number",
				}),
			},
		).Create(&worksDeptUser).Error

		if err != nil {
			return err
		}

		err = tx.Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "user_hash"},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"ko_lang",
					"en_lang",
					"zh_lang",
					"jp_lang",
					"ru_lang",
					"vi_lang",
					"def_lang",
				}),
			},
		).Create(&worksUserMultiLang).Error

		if err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (r *orgRepositoryImpl) InitWorksOrgCode(ctx context.Context) ([]string, error) {

	result := make([]string, 0)
	var model []model.WorksOrgCode

	viewSql := `SELECT org_code FROM org.works_org_code `
	err := r.db.Raw(viewSql).Scan(&model).Error

	if err != nil {
		log.Println("[InitWorksOrgCode] - No record found or DB error")
		return nil, err
	}

	if len(model) > 0 {
		for _, value := range model {
			result = append(result, value.OrgCode)
		}
	}

	return result, nil
}
