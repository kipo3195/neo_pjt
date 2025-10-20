package adapter

import (
	"org/internal/application/usecase/input"
	"org/internal/delivery/dto/department"
)

func CreateDeptInput(req department.CreateDeptRequest) input.CreateDeptInput {
	return input.CreateDeptInput{
		DeptOrg:        req.DeptCode,
		DeptCode:       req.DeptCode,
		ParentDeptCode: req.ParentDeptCode,
		KoLang:         req.KoLang,
		EnLang:         req.EnLang,
		JpLang:         req.JpLang,
		ZhLang:         req.ZhLang,
		RuLang:         req.RuLang,
		ViLang:         req.ViLang,
		Header:         req.Header,
	}
}
