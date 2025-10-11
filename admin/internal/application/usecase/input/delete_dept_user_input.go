package input

type DeleteDeptUserInput struct {
	UserHash string
	DeptCode string
	DeptOrg  string
}

func MakeDeleteDeptUserInput(userHash string, deptCode string, deptOrg string) DeleteDeptUserInput {

	return DeleteDeptUserInput{
		UserHash: userHash,
		DeptCode: deptCode,
		DeptOrg:  deptOrg,
	}
}
