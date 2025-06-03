package dto

type CreateDeptRequest struct {
	DeptCode       string `json:"deptCode" validate:"required"`
	DeptOrg        string `json:"deptOrg" validate:"required"`
	ParentDeptCode string `json:"parentDeptCode" validate:"required"`
	KrLang         string `json:"krLang"`
	EnLang         string `json:"enLang"`
	JpLang         string `json:"jpLang"`
	CnLang         string `json:"cnLang"`
}
