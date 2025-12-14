package repository

import (
	"batch/internal/domain/orgInfo/entity"
	"batch/internal/domain/orgInfo/repository"
	"batch/internal/infrastructure/model"
	"context"
	"log"

	"gorm.io/gorm"
)

type orgInfoRepositoryImpl struct {
	db *gorm.DB
}

func OrgInfoMigrate(db *gorm.DB) {
	db.AutoMigrate(model.OrgInfo{})
	db.AutoMigrate(model.OrgInfoJsonHistory{})
}

func NewOrgInfoRepository(db *gorm.DB) repository.OrgInfoRepository {

	return &orgInfoRepositoryImpl{
		db: db,
	}
}

func (r *orgInfoRepositoryImpl) GetOrgInfo(ctx context.Context, org string) ([]entity.OrgInfoEntity, error) {

	orgInfoModel := []model.OrgInfo{}

	viewSql := `select * from org_info_view where org = ? `
	err := r.db.Raw(viewSql, org).Scan(&orgInfoModel).Error

	if err != nil {
		log.Println("[GetOrgInfo] - No record found or DB error")
		return nil, err
	}

	return convertOrgInfoToEntity(orgInfoModel), nil
}

// 변환 함수
func convertOrgInfoToEntity(models []model.OrgInfo) []entity.OrgInfoEntity {
	entities := make([]entity.OrgInfoEntity, 0, len(models))
	for _, m := range models {
		entity := entity.OrgInfoEntity{
			Org:            m.Org,
			DeptCode:       m.DeptCode,
			ParentDeptCode: m.ParentDeptCode,
			KoLang:         m.KoLang,
			EnLang:         "",
			ZhLang:         "",
			JpLang:         "",
			RuLang:         "",
			ViLang:         "",
			UpdateHash:     "",
			Kind:           m.Kind,
			UserHash:       m.UserHash,
			UserId:         m.UserId,
			Header:         "",
			Description:    "",
		}
		entities = append(entities, entity)
	}
	return entities
}

func (r *orgInfoRepositoryImpl) PutOrgInfoJson(ctx context.Context, org string, fileName string, orgJson string) error {

	if err := r.db.WithContext(ctx).Create(&model.OrgInfoJsonHistory{
		Org:         org,
		FileName:    fileName,
		OrgInfoJson: orgJson,
	}).Error; err != nil {
		log.Println("[OrgInfoJsonHistory] err : ", err)
		return err
	}
	log.Println("[OrgInfoJsonHistory] insert success.")
	return nil
}
