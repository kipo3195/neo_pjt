package handlers

import "admin/usecases"

type RuleHandler struct {
	usecase usecases.RuleUsecase
}

func NewRuleHandler(r usecases.RuleUsecase) *RuleHandler {
	return &RuleHandler{usecase: r}
}
