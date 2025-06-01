package repositories

import (
	clDto "auth/dto/client"
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
	GetAuth(body *clDto.AuthRequest) (*models.AuthInfo, error)
	PutDeviceToken(token *entities.DeviceToken) (bool, error)
	ToDeviceTokenModel(e *entities.DeviceToken) *models.DeviceToken
	GetValidation(header *clDto.LoginRequestHeader) (bool, error)
}

func NewAuthRepository(db *gorm.DB) AuthRepository {

	return &authRepository{db: db}
}

func (r *authRepository) GetAuth(body *clDto.AuthRequest) (*models.AuthInfo, error) {

	var auth *models.AuthInfo

	fmt.Printf("클라이언트가 전달한 id : %s pw : %s \n", body.Id, body.Password)

	err := r.db.Where("id = ?", body.Id).Where("password = ?", body.Password).First(&auth).Error
	// reflect: reflect.Value.SetString using unaddressable value
	// find안에 auth가 포인터 변수가 아닌 값일때 발생할 수 있음. gorm.Find()는 포인터로 받은 값을 채워 넣어야 합니다.

	if err == gorm.ErrRecordNotFound {
		// SQL 실행 중 에러가 발생한 경우
		log.Println("[GetAuth] - No record found or DB error")
		return nil, gorm.ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return auth, err
}

func (r *authRepository) PutDeviceToken(token *entities.DeviceToken) (bool, error) {

	// entity -> model
	models := r.ToDeviceTokenModel(token)

	// Insert 실행
	if err := r.db.Create(&models).Error; err != nil {
		log.Println("[PutDeviceToken] - DB error")
		return false, err
	}
	return true, nil
}

func (r *authRepository) ToDeviceTokenModel(e *entities.DeviceToken) *models.DeviceToken {
	return &models.DeviceToken{
		Uuid:  e.Uuid,
		Token: e.Token,
	}
}

func (r *authRepository) GetValidation(header *clDto.LoginRequestHeader) (bool, error) {

	var validation models.DeviceToken

	fmt.Println("클라이언트가 전달한 토큰 : ", header.Token)

	result := r.db.Where("uuid = ?", header.Uuid).Order("seq DESC").First(&validation)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("[GetValidation] - No record found")
		return false, gorm.ErrRecordNotFound
	} else if result.Error != nil {
		log.Println("[GetValidation] - DB error")
		return false, result.Error
	} else {
		serverToken := validation.Token
		fmt.Printf("UUID %s 로 조회된 토큰 : %s \n", header.Uuid, serverToken)
		if serverToken == header.Token {
			// 토큰 일치
			return true, nil
		} else {
			// 토큰 불일치
			log.Println("[GetValidation] - token mismatch")
			return false, nil
		}
	}
}
