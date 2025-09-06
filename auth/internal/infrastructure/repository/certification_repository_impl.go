package repository

import (
	"auth/internal/domain/certification/entity"
	"auth/internal/domain/certification/repository"
	"auth/internal/infrastructure/model"
	"log"

	"gorm.io/gorm"
)

type cerificationRepository struct {
	db *gorm.DB
}

func NewCertificationRepository(db *gorm.DB) repository.CerificationRepository {
	return &cerificationRepository{db: db}
}

func CertificationMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.AuthInfo{})
}

func (r *cerificationRepository) CheckAuth(en entity.AuthInfoEntity) (*entity.AuthInfoEntity, error) {

	var auth *entity.AuthInfoEntity

	err := r.db.Where("id = ?", en.Id).Where("password = ?", en.Password).First(&auth).Error

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

func (r *cerificationRepository) GetUserHash(en entity.AuthInfoEntity) (string, error) {

	var userHash string

	result := r.db.Raw(`
		SELECT su.user_hash 
		FROM service_users su
		JOIN auth_info ai ON su.user_id = ai.id
		WHERE su.use_yn = 'Y' AND ai.id = ?`, en.Id).Scan(&userHash)

	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected == 0 {
		return "", gorm.ErrRecordNotFound
	}

	return userHash, nil
}

func (r *cerificationRepository) GetValidation(en entity.AppTokenValidationEntity) (bool, error) {

	var validation model.IssuedAppToken

	log.Println("클라이언트가 전달한 토큰 : ", en.AppToken)

	result := r.db.Where("uuid = ?", en.Uuid).Order("seq DESC").First(&validation)

	if result.Error != nil {
		return false, result.Error
	} else {
		serverToken := validation.AppToken
		if serverToken == en.AppToken {
			// 토큰 일치
			return true, nil
		} else {
			// 토큰 불일치
			log.Println("[GetValidation] - token mismatch")
			return false, nil
		}
	}

}
