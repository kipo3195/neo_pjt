package repository

import (
	"auth/internal/consts"
	"auth/internal/domain/shared"
	"auth/internal/domain/token/entity"
	"auth/internal/domain/token/repository"
	"auth/internal/infrastructure/model"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) repository.TokenRepository {
	return &tokenRepository{
		db: db,
	}
}

func TokenMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.IssuedAppToken{})
	db.AutoMigrate(&model.IssuedAuthTokenHistory{})
	db.AutoMigrate(&model.AuthTokenInfo{})
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

func (r *tokenRepository) InitUserAuthToken() ([]entity.AuthTokenEntity, error) {

	var histories []model.IssuedAuthTokenHistory

	// 서브쿼리: id, uuid 기준으로 최신 create_at 찾기
	subQuery := r.db.Model(&model.IssuedAuthTokenHistory{}).
		Select("id, uuid, MAX(create_at) as max_create_at").
		Group("id, uuid")

	// 메인쿼리: 해당 최신 create_at 레코드 가져오기
	err := r.db.Table("issued_auth_token_history h").
		Joins("JOIN (?) m ON h.id = m.id AND h.uuid = m.uuid AND h.create_at = m.max_create_at", subQuery).
		Find(&histories).Error

	if err != nil {
		panic(err)
	}

	return ToAuthTokenEntities(histories), nil
}

func ToAuthTokenEntities(histories []model.IssuedAuthTokenHistory) []entity.AuthTokenEntity {
	result := make([]entity.AuthTokenEntity, len(histories))
	for i, h := range histories {
		result[i] = ToAuthTokenEntity(h)
	}
	return result
}

func ToAuthTokenEntity(h model.IssuedAuthTokenHistory) entity.AuthTokenEntity {
	return entity.AuthTokenEntity{
		Id:   h.Id,
		Uuid: h.Uuid,
		At:   h.AccessToken,
		Rt:   h.RefreshToken,
	}
}

func (r *tokenRepository) InitAuthTokenInfo(ctx context.Context) ([]entity.AuthTokenInfoEntity, error) {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var authTokenInfo []model.AuthTokenInfo
	if err := tx.Find(&authTokenInfo).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// model -> entity
	var entities []entity.AuthTokenInfoEntity
	for _, m := range authTokenInfo {
		e := entity.AuthTokenInfoEntity{
			TokenType: m.TokenType,
			TokenExp:  m.TokenExp,
		}
		entities = append(entities, e)
	}

	// 트랜잭션 종료
	if err := tx.Commit().Error; err != nil {
		log.Println("[InitDeviceTokenInfo] - Commit failed")
		return nil, consts.ErrDB
	}
	log.Println("[InitDeviceTokenInfo] - Commit Success")

	return entities, nil
}

func (r *tokenRepository) PutAuthToken(ctx context.Context, id string, uuid string, at string, rt string, rtExp string) error {
	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(&model.IssuedAuthTokenHistory{
		Id:              id,
		Uuid:            uuid,
		AccessToken:     at,
		RefreshToken:    rt,
		RefreshTokenExp: rtExp,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 트랜잭션 종료
	if err := tx.Commit().Error; err != nil {
		log.Println("[PutAuthToken] - Commit failed")
		return consts.ErrDB
	}
	log.Println("[PutAuthToken] - Commit Success")
	return nil
}
