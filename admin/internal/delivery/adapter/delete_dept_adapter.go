package adapter

import "admin/internal/application/usecase/input"

func MakeDeleteDeptInput(deptOrg string, deptCode string) input.DeleteDeptInput {
	return input.DeleteDeptInput{
		DeptOrg:  deptOrg,
		DeptCode: deptCode,
	}
}
