package repository

import (
	"batch/internal/domain/extendDbConnect/repository"

	"gorm.io/gorm"
)

type extendDbConnectRepositoryImpl struct {
}

func NewExtendDBConnectRepository(db *gorm.DB) repository.ExtendDBConnectRepository {

	return &extendDbConnectRepositoryImpl{}
}
