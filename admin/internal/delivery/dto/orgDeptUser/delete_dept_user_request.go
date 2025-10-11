package orgDeptUser

type DeleteDeptUserRequest struct {
	UserHash string `json:"userHash" validate:"required"`
	DeptCode string `json:"deptCode" validate:"required"`
	DeptOrg  string `json:"deptOrg" validate:"required"`
}
