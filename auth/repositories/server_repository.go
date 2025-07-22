package repositories

import (
	"auth/entities"
	"auth/models"
	"log"

	"gorm.io/gorm"
)

type serverRepository struct {
	db *gorm.DB
}

type ServerRepository interface {
	PutIssuedAppToken(token *entities.AppTokenEntity) (bool, error)
	GetValidation(entity entities.AppTokenValidationEntity) (bool, error)
}

func NewServerRepository(db *gorm.DB) ServerRepository {
	return &serverRepository{
		db: db,
	}
}

func (r *serverRepository) PutIssuedAppToken(token *entities.AppTokenEntity) (bool, error) {

	// entity -> model
	models := toAppTokenModel(token)

	// Insert 실행
	if err := r.db.Create(&models).Error; err != nil {
		log.Println("[PutAppToken] - DB error")
		return false, err
	}

	return true, nil
}

func (r *serverRepository) GetValidation(entity entities.AppTokenValidationEntity) (bool, error) {

	var validation models.IssuedAppToken

	log.Println("클라이언트가 전달한 토큰 : ", entity.AppToken)

	result := r.db.Where("uuid = ?", entity.Uuid).Order("seq DESC").First(&validation)

	if result.Error != nil {
		return false, result.Error
	} else {
		serverToken := validation.AppToken
		if serverToken == entity.AppToken {
			// 토큰 일치
			return true, nil
		} else {
			// 토큰 불일치
			log.Println("[GetValidation] - token mismatch")
			return false, nil
		}
	}

}
