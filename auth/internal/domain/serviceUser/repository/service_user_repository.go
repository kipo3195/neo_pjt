package repository

import "auth/internal/domain/serviceUser/entity"

type serviceUserRepository struct {
}

type ServiceUserRepository interface {
	InitServiceUsers() ([]entity.ServiceUserEntity, error)
}
