package adapter

import "admin/internal/application/usecase/input"

func MakeRegistDeptUserInput(userHash string, deptCode string, deptOrg string, positionCode string, roleCode string, isConcurrentPosition string) input.RegistDeptUserInput {
	return input.RegistDeptUserInput{
		UserHash:             userHash,
		DeptCode:             deptCode,
		DeptOrg:              deptOrg,
		PositionCode:         positionCode,
		RoleCode:             roleCode,
		IsConcurrentPosition: isConcurrentPosition,
	}
}
