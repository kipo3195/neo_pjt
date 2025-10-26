package adapter

import "admin/internal/application/usecase/input"

func MakeDeleteDeptUserInput(userHash string, deptCode string, deptOrg string) input.DeleteDeptUserInput {

	return input.DeleteDeptUserInput{
		UserHash: userHash,
		DeptCode: deptCode,
		DeptOrg:  deptOrg,
	}
}
