package server

type ServerCreateDeptRequest struct {
	DeptOrg        string `json:"deptOrg"`
	DeptCode       string `json:"deptCode"`
	ParentDeptCode string `json:"parentDeptCode"`
	DeptNameKr     string `json:"deptNameKr"`
	DeptNameEn     string `json:"deptNameEn"`
	DeptNameJp     string `json:"deptNameJp"`
	DeptNameCn     string `json:"deptNameCn"`
}
