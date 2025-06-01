package dto

type CreateDeptRequest struct {
	DeptOrg        string `json:"deptOrg" validate:"required"`
	DeptCode       string `json:"deptCode" validate:"required"`
	ParentDeptCode string `json:"parentDeptCode" validate:"required"`
	KrLang         string `json:"krLang"`
	EnLang         string `json:"enLang"`
	JpLang         string `json:"jpLang"`
	CnLang         string `json:"cnLang"`
}
