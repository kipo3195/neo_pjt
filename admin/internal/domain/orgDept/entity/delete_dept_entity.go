package entity

type DeleteDeptEntity struct {
	DeptOrg  string `json:"deptOrg"`
	DeptCode string `json:"deptCode"`
}

func MakeDeleteDeptEntity(deptOrg string, deptCode string) DeleteDeptEntity {
	return DeleteDeptEntity{
		DeptOrg:  deptOrg,
		DeptCode: deptCode,
	}
}
