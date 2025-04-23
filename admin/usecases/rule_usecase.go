package usecases

import (
	"admin/repositories"
)

type ruleUsecase struct {
	repo repositories.RuleRepository
}

type RuleUsecase interface {
}

func NewRuleUsecase(repo repositories.RuleRepository) RuleUsecase {
	return &ruleUsecase{repo: repo}
}
