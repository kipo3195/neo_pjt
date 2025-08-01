package repositories

import (
	"auth/internal/domains/certification/entities"
	"auth/internal/domains/certification/models"
	sharedModels "auth/internal/sharedModels"

	"log"

	"gorm.io/gorm"
)

type cerificationRepository struct {
	db *gorm.DB
}

type CerificationRepository interface {
	CheckAuth(entity entities.AuthInfoEntity) (*models.AuthInfo, error)
	GetUserHash(entity entities.AuthInfoEntity) (string, error)
	GetValidation(entity entities.AppTokenValidationEntity) (bool, error)
}

func NewCertificationRepository(db *gorm.DB) CerificationRepository {
	return &cerificationRepository{db: db}
}

func (r *cerificationRepository) CheckAuth(entity entities.AuthInfoEntity) (*models.AuthInfo, error) {

	var auth *models.AuthInfo

	err := r.db.Where("id = ?", entity.Id).Where("password = ?", entity.Password).First(&auth).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// SQL 실행 중 에러가 발생한 경우
			log.Println("[GetAuth] - No record found or DB error")
			return nil, gorm.ErrRecordNotFound
		} else {
			return nil, err
		}
	}
	return auth, nil
}

func (r *cerificationRepository) GetUserHash(entity entities.AuthInfoEntity) (string, error) {

	var userHash string

	result := r.db.Raw(`
		SELECT su.user_hash 
		FROM service_users su
		JOIN auth_info ai ON su.user_id = ai.id
		WHERE su.use_yn = 'Y' AND ai.id = ?`, entity.Id).Scan(&userHash)

	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected == 0 {
		return "", gorm.ErrRecordNotFound
	}

	return userHash, nil
}

// func toAppTokenModel(e *pkgEntities.AppTokenEntity) *pkgModels.IssuedAppToken {
// 	return &pkgModels.IssuedAppToken{
// 		Uuid:         e.Uuid,
// 		AppToken:     e.AppToken,
// 		RefreshToken: e.RefreshToken,
// 	}
// }

func (r *cerificationRepository) GetValidation(entity entities.AppTokenValidationEntity) (bool, error) {

	var validation sharedModels.IssuedAppToken

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
