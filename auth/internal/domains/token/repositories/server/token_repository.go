package repositories

import (
	"auth/internal/domains/token/entities"
	sharedEntities "auth/internal/sharedEntities"
	sharedModels "auth/internal/sharedModels"
	"errors"
	"log"

	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

type TokenRepository interface {
	PutIssuedAppToken(token *sharedEntities.AppTokenEntity) (bool, error)
	GetValidationAppToken(entity entities.AppTokenValidationEntity) (bool, error)
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{
		db: db,
	}
}

func (r *tokenRepository) PutIssuedAppToken(token *sharedEntities.AppTokenEntity) (bool, error) {

	// entity -> model
	issuedAppToken := toAppTokenModel(token)

	// Insert 실행
	if err := r.db.Create(&issuedAppToken).Error; err != nil {
		log.Println("[PutAppToken] - DB error")
		return false, err
	}

	return true, nil
}

func (r *tokenRepository) GetValidationAppToken(entity entities.AppTokenValidationEntity) (bool, error) {

	var validation sharedModels.IssuedAppToken

	log.Println("클라이언트가 전달한 토큰 : ", entity.Token)

	result := r.db.Where("uuid = ?", entity.Uuid).Order("seq DESC").First(&validation)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("[GetValidation] result record = 0")
		return false, result.Error
	} else if result.Error != nil {
		log.Println("[GetValidation] DB error")
		return false, result.Error
	} else {
		serverToken := validation.AppToken
		if serverToken == entity.Token {
			// 토큰 일치
			return true, nil
		} else {
			// 토큰 불일치
			log.Println("[GetValidation] - token mismatch")
			return false, nil
		}
	}

}

func toAppTokenModel(e *sharedEntities.AppTokenEntity) *sharedModels.IssuedAppToken {
	return &sharedModels.IssuedAppToken{
		Uuid:         e.Uuid,
		AppToken:     e.AppToken,
		RefreshToken: e.RefreshToken,
	}
}
