package org

type CreateDeptRequest struct {
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
