package server

type ServerDeleteDeptRequest struct {
	DeptOrg  string `json:"deptOrg"`
	DeptCode string `json:"deptCode"`
}
