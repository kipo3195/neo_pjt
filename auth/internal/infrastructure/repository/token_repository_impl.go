package repository

import (
	"auth/internal/domain/shared"
	"auth/internal/domain/token/entity"
	"auth/internal/infrastructure/model"
	"errors"
	"log"

	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

type TokenRepository interface {
	PutIssuedAppToken(token *shared.AppTokenEntity) (bool, error)
	GetValidationAppToken(entity entity.AppTokenValidationEntity) (bool, error)
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{
		db: db,
	}
}

func TokenMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.IssuedAppToken{})
}

func (r *tokenRepository) PutIssuedAppToken(token *shared.AppTokenEntity) (bool, error) {

	// entity -> model
	issuedAppToken := toAppTokenModel(token)

	// Insert 실행
	if err := r.db.Create(&issuedAppToken).Error; err != nil {
		log.Println("[PutAppToken] - DB error")
		return false, err
	}

	return true, nil
}

func (r *tokenRepository) GetValidationAppToken(entity entity.AppTokenValidationEntity) (bool, error) {

	var validation model.IssuedAppToken

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

func toAppTokenModel(e *shared.AppTokenEntity) *model.IssuedAppToken {
	return &model.IssuedAppToken{
		Uuid:         e.Uuid,
		AppToken:     e.AppToken,
		RefreshToken: e.RefreshToken,
	}
}
