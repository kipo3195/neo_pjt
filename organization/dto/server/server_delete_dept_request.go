package server

type SvDeleteDeptRequest struct {
	DeptOrg  string `json:"deptOrg"`
	DeptCode string `json:"deptCode"`
}
