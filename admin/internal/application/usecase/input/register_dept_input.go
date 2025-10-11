package input

import "admin/internal/delivery/dto/orgDept"

type RegisterDeptInput struct {
	DeptCode       string `json:"deptCode" validate:"required"`
	DeptOrg        string `json:"deptOrg" validate:"required"`
	ParentDeptCode string `json:"parentDeptCode" validate:"required"`
	KoLang         string `json:"ko"`
	EnLang         string `json:"en"`
	JpLang         string `json:"jp"`
	ZhLang         string `json:"zh"`
	RuLang         string `json:"ru"`
	ViLang         string `json:"vi"`
	Header         string `json:"header"` // 부서장
}

func MakeRegisterDeptInput(req orgDept.RegisterDeptRequest) RegisterDeptInput {
	return RegisterDeptInput{
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
