package repositories

import (
	customErr "auth/consts"
	"auth/entities"
	"auth/models"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

type AuthRepository interface {
	CheckAuth(entity entities.AuthInfoEntity) (*models.AuthInfo, error)
	GetUserHash(entity entities.AuthInfoEntity) (string, error)
	PutDeviceToken(token *entities.DeviceTokenEntity) (bool, error)
	ToDeviceTokenModel(e *entities.DeviceTokenEntity) *models.DeviceToken
	GetValidation(entity entities.AppTokenValidationEntity) (bool, error)
}

func NewAuthRepository(db *gorm.DB) AuthRepository {

	return &authRepository{db: db}
}

func (r *authRepository) CheckAuth(entity entities.AuthInfoEntity) (*models.AuthInfo, error) {

	var auth *models.AuthInfo

	err := r.db.Where("id = ?", entity.Id).Where("password = ?", entity.Password).First(&auth).Error

	if err == gorm.ErrRecordNotFound {
		// SQL 실행 중 에러가 발생한 경우
		log.Println("[GetAuth] - No record found or DB error")
		return nil, gorm.ErrRecordNotFound
	}
	return auth, nil
}

func (r *authRepository) GetUserHash(entity entities.AuthInfoEntity) (string, error) {
	var user struct {
		UserHash string
	}

	result := r.db.Raw(`
		SELECT su.user_hash 
		FROM service_users su
		JOIN auth_info ai ON su.user_id = ai.id
		WHERE su.use_yn = 'Y' AND ai.id = ?`, entity.Id).Scan(&user)

	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected == 0 {
		return "", gorm.ErrRecordNotFound
	}

	return user.UserHash, nil
}

func (r *authRepository) PutDeviceToken(token *entities.DeviceTokenEntity) (bool, error) {

	// entity -> model
	models := r.ToDeviceTokenModel(token)

	// Insert 실행
	if err := r.db.Create(&models).Error; err != nil {
		log.Println("[PutDeviceToken] - DB error")
		return false, err
	}
	return true, nil
}

func (r *authRepository) ToDeviceTokenModel(e *entities.DeviceTokenEntity) *models.DeviceToken {
	return &models.DeviceToken{
		Uuid:  e.Uuid,
		Token: e.Token,
	}
}

func (r *authRepository) GetValidation(entity entities.AppTokenValidationEntity) (bool, error) {

	var validation models.DeviceToken

	fmt.Println("클라이언트가 전달한 토큰 : ", entity.AppToken)

	result := r.db.Where("uuid = ?", entity.Uuid).Order("seq DESC").First(&validation)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("[GetValidation] - No record found")
		return false, customErr.ErrDbRowNotFound
	} else if result.Error != nil {
		log.Println("[GetValidation] - DB error")
		return false, customErr.ErrDB
	} else {
		serverToken := validation.Token
		fmt.Printf("UUID %s 로 조회된 토큰 : %s \n", entity.Uuid, serverToken)
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
