package repositories

import (
	"context"
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
}

func NewOrgRepository(db *gorm.DB) OrgRepository {
	return &orgRepository{db: db}
}

func (r *orgRepository) SaveDepartment(ctx context.Context, entity entities.CreateDepartmentEntity) (interface{}, error) {

	models := r.ToDepartmentModel(entity)

	if err := r.db.Create(&models).Error; err != nil {
		log.Println("[SaveDepartment] - DB error")
		return false, err
	}
	return true, nil
}

func (r *orgRepository) ToDepartmentModel(e entities.CreateDepartmentEntity) models.Department {
	return models.Department{
		DeptName: e.DeptName,
	}
}
