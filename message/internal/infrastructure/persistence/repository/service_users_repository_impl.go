package repository

import (
	"message/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

// 도메인 로직이 아닌 단순 테이블 생성 용
func ServiceUsersMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.ServiceUsers{})
}
