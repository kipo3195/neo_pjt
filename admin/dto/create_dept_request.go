package dto

type CreateDeptRequest struct {
	DeptOrg        string `json:"deptOrg" validate:"required"`
	DeptCode       string `json:"deptCode" validate:"required"`
	ParentDeptCode string `json:"parentDeptCode" validate:"required"`
	DeptNameKr     string `json:"deptNameKr"`
	DeptNameEn     string `json:"deptNameEn"`
	DeptNameJp     string `json:"deptNameJp"`
	DeptNameCn     string `json:"deptNameCn"`
}
