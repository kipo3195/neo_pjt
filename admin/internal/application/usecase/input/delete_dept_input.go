package input

type DeleteDeptInput struct {
	DeptOrg  string
	DeptCode string
}

func MakeDeleteDeptInput(deptOrg string, deptCode string) DeleteDeptInput {
	return DeleteDeptInput{
		DeptOrg:  deptOrg,
		DeptCode: deptCode,
	}
}
