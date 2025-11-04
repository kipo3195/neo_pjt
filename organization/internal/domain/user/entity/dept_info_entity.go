package entity

type DeptInfoEntity struct {
	DeptOrg      string             `json:"deptOrg"`
	DeptCode     string             `json:"deptCode"`
	DeptName     DeptNameEntity     `json:"deptName"`
	Header       string             `json:"header"`
	Description  string             `json:"description"`
	RoleName     RoleNameEntity     `json:"roleName"`
	PositionName PositionNameEntity `json:"positionName"`
}
