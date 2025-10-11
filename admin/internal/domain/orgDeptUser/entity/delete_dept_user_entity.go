package entity

type DeleteDeptUserEntity struct {
	UserHash string `json:"userHash" validate:"required"`
	DeptCode string `json:"deptCode" validate:"required"`
	DeptOrg  string `json:"deptOrg" validate:"required"`
}

func MakeDeleteDeptUserEntity(userhash string, deptOrg string, deptCode string) DeleteDeptUserEntity {

	return DeleteDeptUserEntity{
		UserHash: userhash,
		DeptCode: deptCode,
		DeptOrg:  deptOrg,
	}
}
