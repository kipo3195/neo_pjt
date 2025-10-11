package input

type RegistDeptUserInput struct {
	UserHash             string
	DeptCode             string
	DeptOrg              string
	PositionCode         string
	RoleCode             string
	IsConcurrentPosition string
}

func MakeRegistDeptUserInput(userHash string, deptCode string, deptOrg string, positionCode string, roleCode string, isConcurrentPosition string) RegistDeptUserInput {
	return RegistDeptUserInput{
		UserHash:             userHash,
		DeptCode:             deptCode,
		DeptOrg:              deptOrg,
		PositionCode:         positionCode,
		RoleCode:             roleCode,
		IsConcurrentPosition: isConcurrentPosition,
	}
}
