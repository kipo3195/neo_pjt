package repositories

import (
	"gorm.io/gorm"
)

type ruleRepository struct {
	db *gorm.DB
}

type RuleRepository interface {
}

func NewRuleRepository(db *gorm.DB) RuleRepository {
	return &ruleRepository{db: db}
}
