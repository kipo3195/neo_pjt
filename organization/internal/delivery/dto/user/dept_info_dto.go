package user

type DeptInfoDto struct {
	DeptOrg      string          `json:"deptOrg"`
	DeptCode     string          `json:"deptCode"`
	DeptName     DeptNameDto     `json:"deptName"`
	Header       string          `json:"header"`
	Description  string          `json:"description"`
	PositionName PositionNameDto `json:"positionName"`
	RoleName     RoleNameDto     `json:"roleName"`
}
