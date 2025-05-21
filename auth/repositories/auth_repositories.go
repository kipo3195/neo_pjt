package repositories

import (
	"auth/dto"
	"auth/entities"
	"auth/models"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

type AuthRepository interface {
	GetAuth(dto.AuthRequest) (*models.AuthInfo, error)
	PutDeviceToken(token *entities.DeviceToken) (bool, error)
	ToEntityDeviceToken(m *entities.DeviceToken) *models.DeviceToken
}

func NewAuthRepository(db *gorm.DB) AuthRepository {

	return &authRepository{db: db}
}

func (r *authRepository) GetAuth(req dto.AuthRequest) (*models.AuthInfo, error) {

	var auth models.AuthInfo

	err := r.db.Where("id = ?", req.Id).Where("password = ?", req.Password).Find(&auth).Error
	// reflect: reflect.Value.SetString using unaddressable value
	// find안에 auth가 포인터 변수가 아닌 값일때 발생할 수 있음. gorm.Find()는 포인터로 받은 값을 채워 넣어야 합니다.

	return &auth, err
}

func (r *authRepository) PutDeviceToken(token *entities.DeviceToken) (bool, error) {

	// entity -> model
	models := r.ToEntityDeviceToken(token)

	// Insert 실행
	if err := r.db.Create(&models).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *authRepository) ToEntityDeviceToken(m *entities.DeviceToken) *models.DeviceToken {
	return &models.DeviceToken{
		Uuid:  m.Uuid,
		Token: m.Token,
	}
}
