package adapter

import (
	"admin/internal/application/usecase/input"
	"admin/internal/delivery/dto/orgDept"
)

func MakeRegisterDeptInput(req orgDept.RegisterDeptRequest) input.RegisterDeptInput {
	return input.RegisterDeptInput{
		DeptCode:       req.DeptCode,
		DeptOrg:        req.DeptOrg,
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
